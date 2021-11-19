package controllers

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

func TestGenericResponse_Deserialize(t *testing.T) {
	type fields struct {
		Success bool
		Error   string
	}
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantFields fields
	}{
		{
			name: "deserializes",
			args: args{
				in: bytes.NewBufferString(`{"success": true, "error":"error"}`),
			},
			wantErr: false,
			wantFields: fields{
				Success: true,
				Error:   "error",
			},
		},
		{
			name: "fails",
			args: args{
				in: bytes.NewBufferString(`{rietnsriten :"asretnerror"}`),
			},
			wantErr: true,
			wantFields: fields{
				Success: false,
				Error:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GenericResponse{}
			if err := gr.Deserialize(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf(
					"GenericResponse.Deserialize() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if gr.Success != tt.wantFields.Success {
				t.Errorf(
					"GenericResponse.Deserialize() gr = %v, want %v",
					gr,
					tt.wantFields,
				)
			}
			if gr.Error != tt.wantFields.Error {
				t.Errorf(
					"GenericResponse.Deserialize() gr = %v, want %v",
					gr,
					tt.wantFields,
				)
			}
		})
	}
}

func TestGenericResponse_Serialize(t *testing.T) {
	type fields struct {
		Success bool
		Error   string
	}
	type args struct {
		out *string
	}
	okinput := ""
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantOut string
	}{
		{
			name: "serializes",
			fields: fields{
				Success: true,
				Error:   "",
			},
			args: args{
				out: &okinput,
			},
			wantErr: false,
			wantOut: "{\"success\":true,\"error\":\"\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GenericResponse{
				Success: tt.fields.Success,
				Error:   tt.fields.Error,
			}
			if err := gr.Serialize(tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf(
					"GenericResponse.Serialize() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if *tt.args.out != tt.wantOut {
				t.Errorf(
					"GenericResponse.Serialize() out = %v, wantOut %v",
					*tt.args.out,
					tt.wantOut,
				)
			}
		})
	}
}

func findEnv(project models.Project, name string) *models.Environment {
	for _, env := range project.Environments {
		if env.Name == name {
			return &env
		}
	}

	return nil
}

func findMessages(
	messages []models.Message,
	user models.User,
	environments []models.Environment,
) map[string]models.Message {
	m := make(map[string]models.Message)

	for _, message := range messages {
		for _, environment := range environments {
			if message.RecipientID == user.ID &&
				message.EnvironmentID == environment.EnvironmentID {
				m[environment.Name] = message
			}
		}
	}

	return m
}

func TestGetMessagesFromProjectByUser(t *testing.T) {
	project, users, messages := seedMessages()
	defer teardownMessages(project, users, messages)

	devEnvironment := findEnv(project, "dev")
	stagingEnvironment := findEnv(project, "staging")
	prodEnvironment := findEnv(project, "prod")

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name             string
		args             args
		want             *models.GetMessageByEnvironmentResponse
		wantEnvironments []string
		wantStatus       int
		wantErr          string
	}{
		{
			name: "gets message for a user - dev env for dev user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["developer"],
			},
			wantEnvironments: []string{"dev"},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{
					"dev": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *devEnvironment,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "gets message for a user - dev env for lead-dev user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["lead-dev"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["lead-dev"],
			},
			wantEnvironments: []string{"dev"},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{
					"dev": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *devEnvironment,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "gets message for a user - all envs for devops user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["devops"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["devops"],
			},
			wantEnvironments: []string{"dev", "staging", "prod"},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{
					"dev": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *devEnvironment,
					},
					"staging": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *stagingEnvironment,
					},
					"prod": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *prodEnvironment,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "gets message for a user - all envs for admit user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["admin"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["admin"],
			},
			wantEnvironments: []string{"dev", "staging", "prod"},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{
					"dev": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *devEnvironment,
					},
					"staging": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *stagingEnvironment,
					},
					"prod": {
						Message: findMessages(
							messages,
							users["developer"],
							project.Environments,
						)["dev"],
						Environment: *prodEnvironment,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "project not found",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "not a project",
					"device":    users["admin"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["admin"],
			},
			wantEnvironments: []string{},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "device not found",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    "not a device",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["admin"],
			},
			wantEnvironments: []string{},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "wrong user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["admin"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: users["dev"],
			},
			wantEnvironments: []string{},
			want: &models.GetMessageByEnvironmentResponse{
				Environments: map[string]models.GetMessageResponse{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetMessagesFromProjectByUser(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetMessagesFromProjectByUser() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetMessagesFromProjectByUser() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.(*models.GetMessageByEnvironmentResponse)
			for _, wantE := range tt.wantEnvironments {
				if gotE, ok := gotResponse.Environments[wantE]; ok {
					gotEnvironment := gotE.Environment
					wantEnvironment := tt.want.Environments[wantE].Environment

					gotMessage := gotE.Message
					wantMessage := tt.want.Environments[wantE].Message

					if gotEnvironment.EnvironmentID != wantEnvironment.EnvironmentID {
						t.Errorf(
							"GetMessagesFromProjectByUser() Environments[%v].Environment = %v, want %v",
							wantE,
							gotEnvironment,
							wantEnvironment,
						)
					}
					if gotMessage.Uuid != wantMessage.Uuid {
						t.Errorf(
							"GetMessagesFromProjectByUser() Environments[%v].Message = %v, want %v",
							wantE,
							gotMessage,
							wantMessage,
						)
					}

				} else {
					t.Errorf("GetMessagesFromProjectByUser() Environments[%v] = %v, wantErr %v", wantE, ok, wantE)
				}
			}
		})
	}
}

func TestWriteMessages(t *testing.T) {
	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       router.Serde
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := WriteMessages(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"WriteMessages() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WriteMessages() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"WriteMessages() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       router.Serde
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := DeleteMessage(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"DeleteMessage() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteMessage() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"DeleteMessage() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestDeleteExpiredMessages(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteExpiredMessages(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestAlertMessagesWillExpire(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AlertMessagesWillExpire(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func seedMessages() (
	project models.Project,
	users map[string]models.User,
	messages []models.Message,
) {
	users = make(map[string]models.User)

	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		faker.FakeData(&project)
		db.Save(&project)

		roles := make([]models.Role, 0)
		db.Find(&roles)

		environmentTypes := make([]models.EnvironmentType, 3)
		db.Find(&environmentTypes)

		var devEnvironment models.Environment
		for _, e := range environmentTypes {
			env := models.Environment{
				Name:              e.Name,
				EnvironmentTypeID: e.ID,
				EnvironmentType:   e,
				ProjectID:         project.ID,
				Project:           project,
				VersionID:         faker.UUIDHyphenated(),
				EnvironmentID:     faker.UUIDHyphenated(),
			}
			db.Save(&env)

			switch e.Name {
			case "dev":
				devEnvironment = env
			}
		}

		db.Model(&project).
			Association("Environments").
			Find(&project.Environments)

		for _, role := range roles {
			user := models.User{}
			faker.FakeData(&user)
			db.Save(&user)

			device := models.Device{}
			faker.FakeData(&device)
			db.Save(&device)

			db.Model(&user).Association("Devices").Append(&device)
			user.Devices = []models.Device{device}

			projectMember := models.ProjectMember{
				UserID:    user.ID,
				ProjectID: project.ID,
				RoleID:    role.ID,
			}

			db.Save(&projectMember)
			users[role.Name] = user
		}

		// Simulate the admin user having sent a message on the
		// developement environment
		adminUser := users["admin"]

		for _, user := range users {
			message := models.Message{
				Payload:           []byte{},
				Sender:            adminUser,
				SenderID:          adminUser.ID,
				Recipient:         user,
				RecipientID:       user.ID,
				Environment:       devEnvironment,
				EnvironmentID:     devEnvironment.EnvironmentID,
				Uuid:              faker.UUIDHyphenated(),
				RecipientDeviceID: user.Devices[0].ID,
				SenderDeviceID:    adminUser.Devices[0].ID,
			}

			db.Save(&message)

			messages = append(messages, message)
		}

		return nil
	})

	return project, users, messages
}

func teardownMessages(
	project models.Project,
	users map[string]models.User,
	messages []models.Message,
) {
	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(messages)
		db.Delete(project.Environments)

		for _, user := range users {
			db.Model(&user).Association("Devices").Delete(user.Devices)
			db.Delete(
				&models.ProjectMember{},
				"user_id = ? and project_id = ?",
				user.ID,
				project.ID,
			)
			db.Delete(&user)
		}

		organization := project.Organization
		db.Delete(&project)
		db.Delete(&organization)

		return nil
	})
}
