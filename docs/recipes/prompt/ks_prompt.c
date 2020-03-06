#include <dirent.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#if defined(WIN32) || defined(_WIN32)
#define PATH_SEPARATOR "\\"
#else
#define PATH_SEPARATOR "/"
#endif

static char *KS_CONFIG_FILE = "/.ksconfig";

int split(const char *str, char c, char ***arr) {
  int count = 1;
  int token_len = 1;
  int i = 0;
  char *p;
  char *t;

  p = str;
  while (*p != '\0') {
    if (*p == c)
      count++;
    p++;
  }

  *arr = (char **)malloc(sizeof(char *) * count);
  if (*arr == NULL)
    exit(1);

  p = str;
  while (*p != '\0') {
    if (*p == c) {
      (*arr)[i] = (char *)malloc(sizeof(char) * token_len);
      if ((*arr)[i] == NULL)
        exit(1);

      token_len = 0;
      i++;
    }
    p++;
    token_len++;
  }
  (*arr)[i] = (char *)malloc(sizeof(char) * token_len);
  if ((*arr)[i] == NULL)
    exit(1);

  i = 0;
  p = str;
  t = ((*arr)[i]);
  while (*p != '\0') {
    if (*p != c && *p != '\0') {
      *t = *p;
      t++;
    } else {
      *t = '\0';
      i++;
      t = ((*arr)[i]);
    }
    p++;
  }

  return count;
}

char *ReadFile(char *filename) {
  char *buffer = NULL;
  int string_size, read_size;
  FILE *handler = fopen(filename, "r");

  if (handler) {
    fseek(handler, 0, SEEK_END);
    string_size = ftell(handler);
    rewind(handler);

    buffer = (char *)malloc(sizeof(char) * (string_size + 1));

    read_size = fread(buffer, sizeof(char), string_size, handler);

    buffer[string_size] = '\0';

    if (string_size != read_size) {
      free(buffer);
      buffer = NULL;
    }

    fclose(handler);
  }

  return buffer;
}

int exists(const char *fname) {
  FILE *file;
  if ((file = fopen(fname, "r"))) {
    fclose(file);
    return 1;
  }
  return 0;
}

char *get_ks_root() {
  char *potential_path;
  potential_path = (char *)malloc(100);
  strcpy(potential_path, getcwd(NULL, 0));
  while (strcmp(potential_path, "")) {
    int i;
    char *s, *tofree;
    tofree = s = strdup(potential_path);
    int size = 0;
    char **arr = NULL;

    size = split(potential_path, '/', &arr);

    free(tofree);

    char *parent_dir;
    parent_dir = (char *)malloc(500);
    strcpy(parent_dir, "");
    for (i = 1; i < size - 1; i++) {
      strcat(parent_dir, "/");
      strcat(parent_dir, arr[i]);
    }

    char *path_to_env_file;
    path_to_env_file = (char *)malloc(500);

    strcpy(path_to_env_file, potential_path);
    strcat(path_to_env_file, "/.keystone/envconfig");

    if (exists(path_to_env_file)) {
      return potential_path;
    }
    free(path_to_env_file);

    potential_path = parent_dir;
  }
}

void *print_working_env(int with_k) {
  char *potential_path;
  potential_path = (char *)malloc(100);
  strcpy(potential_path, getcwd(NULL, 0));
  while (strcmp(potential_path, "")) {
    int i;
    char *s, *tofree;
    tofree = s = strdup(potential_path);
    int size = 0;
    char **arr = NULL;

    size = split(potential_path, '/', &arr);

    free(tofree);

    char *parent_dir;
    parent_dir = (char *)malloc(500);
    strcpy(parent_dir, "");
    for (i = 1; i < size - 1; i++) {
      strcat(parent_dir, "/");
      strcat(parent_dir, arr[i]);
    }

    char *path_to_env_file;
    path_to_env_file = (char *)malloc(500);

    strcpy(path_to_env_file, potential_path);
    strcat(path_to_env_file, "/.keystone/envconfig");

    if (exists(path_to_env_file)) {
      char *string = (char *)malloc(200);

      string = ReadFile(path_to_env_file);
      if (string) {
        int i2;
        char *s2, *tofree2;
        tofree2 = s2 = strdup(potential_path);
        int size2 = 0;
        char **arr2 = NULL;

        size2 = split(string, '"', &arr2);
        if (with_k)
          printf("Ꝅ %s", arr2[3]);
        else
          printf("%s", arr2[3]);
        free(string);
      }
      break;
    }
    free(path_to_env_file);

    potential_path = parent_dir;
  }
}

int is_directory(const char *path) {
  struct stat statbuf;
  if (stat(path, &statbuf) != 0)
    return 0;
  return S_ISDIR(statbuf.st_mode);
}

char *replace_str(char *str, char *orig, char *rep) {
  static char buffer[4096];
  char *p;

  if (!(p = strstr(str, orig))) // Is 'orig' even in 'str'?
    return str;

  strncpy(buffer, str,
          p - str); // Copy characters from 'str' start to 'orig' st$
  buffer[p - str] = '\0';

  sprintf(buffer + (p - str), "%s%s", rep, p + strlen(orig));

  return buffer;
}

char *list_files_in_dir(char *dir, char *files_list) {
  struct dirent *de;

  DIR *dr = opendir(dir);
  if (dr == NULL)
    return 0;

  while ((de = readdir(dr)) != NULL) {
    if (strcmp(de->d_name, ".") == 0 || strcmp(de->d_name, "..") == 0)
      continue;

    char *sub_dir;
    sub_dir = (char *)malloc(sizeof(char *) * 50);
    strcpy(sub_dir, dir);
    strcat(sub_dir, PATH_SEPARATOR);
    strcat(sub_dir, de->d_name);

    if (exists(sub_dir) && !is_directory(sub_dir)) {
      strcat(files_list, ";");
      strcat(files_list, sub_dir);
    }

    list_files_in_dir(sub_dir, files_list);
  }
  return files_list;

  closedir(dr);
}

int compare_with_current_changes(char *cache_files_list) {
  int i;
  char *s, *tofree;
  tofree = s = strdup(cache_files_list);
  int size = 0;
  char **arr = NULL;

  size = split(cache_files_list, ';', &arr);

  free(tofree);

  char *parent_dir;
  parent_dir = (char *)malloc(500);
  strcpy(parent_dir, "");
  for (i = 1; i < size; i++) {
    char *cache_content;
    char *current_content;
    cache_content = (char *)malloc(sizeof(char *) * 1000000);
    current_content = (char *)malloc(sizeof(char *) * 1000000);

    cache_content = ReadFile(arr[i]);
    current_content = ReadFile(replace_str(arr[i], ".keystone/cache/", ""));
    if (strcmp(current_content, cache_content))
      return 1;
  }
  return 0;
}

int files_has_changed() {
  char *ks_root_dir;
  ks_root_dir = (char *)malloc(sizeof(char *) * 100);
  ks_root_dir = get_ks_root();

  char *cache_root_dir;
  cache_root_dir = (char *)malloc(sizeof(char *) * 100);
  strcat(cache_root_dir, ks_root_dir);
  strcat(cache_root_dir, "/.keystone/cache");

  char *empty_list;
  empty_list = (char *)malloc(sizeof(char *) * 1000);

  char *cache_files_list;
  cache_files_list = (char *)malloc(sizeof(char *) * 1000);
  cache_files_list = list_files_in_dir(cache_root_dir, empty_list);

  return compare_with_current_changes(cache_files_list);
}

int main(int argc, char **argv) {
  if (argv[1]) {
    if (!strcmp("env", argv[1]))
      print_working_env(0);
    else if (!strcmp("full", argv[1])) {
      print_working_env(1);
      if (files_has_changed())
        puts(" ✘");
      else
        puts(" ✔");
    } else if (!strcmp("status", argv[1])) {
      if (files_has_changed())
        puts("✘");
      else
        puts("✔");
    } else
      return 0;

    return 0;
  }
  return 1;
}
