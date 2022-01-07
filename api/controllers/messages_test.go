package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/message"
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
	project, users, messages := seedMessages(true)
	defer teardownMessages(project, users, messages)

	devEnvironment := findEnv(project, "dev")
	stagingEnvironment := findEnv(project, "staging")
	prodEnvironment := findEnv(project, "prod")

	crashingProject, cusers, cmessages := seedMessages(true)
	defer teardownMessages(crashingProject, cusers, cmessages)

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
				Repo: newFakeRepo(),
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
			name:             "crashes while fetching project",
			args:             args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "crash it",
					"device":    cusers["developer"].Devices[0].UID,
				}),
				in1:    nil,
				Repo:   newFakeRepo(),
				user:   cusers["developer"],
			},
			want:             &models.GetMessageByEnvironmentResponse{},
			wantEnvironments: []string{},
			wantStatus:       http.StatusBadRequest,
			wantErr:          "bad request: unexpected error",
		},
		{
			name: "gets message for a user - dev env for lead-dev user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["lead-dev"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: newFakeRepo(),
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
				Repo: newFakeRepo(),
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
				Repo: newFakeRepo(),
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
				Repo: newFakeRepo(),
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
				Repo: newFakeRepo(),
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
				Repo: newFakeRepo(),
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
	project, users, messages := seedMessages(true)
	defer teardownMessages(project, users, messages)

	devEnvironment := findEnv(project, "dev")
	// stagingEnvironment := findEnv(project, "staging")
	prodEnvironment := findEnv(project, "prod")

	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetEnvironmentsResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "writes a message",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "messages": [
                        {
                            "payload": "PGVuY3J5cHRlZF9jb250ZW50Pg==",
                            "sender_device_uid": "%s",
                            "recipient_device_id": %d,
                            "userid": "%s",
                            "recipient_id": %d,
                            "environment_id": "%s",
                            "update_environment_version": true
                        }
                    ]
                }
                `,
					users["admin"].Devices[0].UID,
					users["devops"].Devices[0].ID,
					users["admin"].UserID,
					users["devops"].ID,
					devEnvironment.EnvironmentID,
				))),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want: &models.GetEnvironmentsResponse{
				Environments: []models.Environment{
					*devEnvironment,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "bad request",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(
					bytes.NewBufferString("not serializable"),
				),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad request: invalid character 'o' in literal null (expecting 'u')",
		},
		{
			name: "empty payload",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "messages": [
                        {
                            "payload": "",
                            "sender_device_uid": "%s",
                            "recipient_device_id": %d,
                            "userid": "%s",
                            "recipient_id": %d,
                            "environment_id": "%s",
                            "update_environment_version": true
                        }
                    ]
                }
                `,
					users["admin"].Devices[0].UID,
					users["devops"].Devices[0].ID,
					users["admin"].UserID,
					users["devops"].ID,
					devEnvironment.EnvironmentID,
				))),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "empty payload cannot be written",
		},
		{
			name: "no such environment",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "messages": [
                        {
                            "payload": "PGVuY3J5cHRlZF9jb250ZW50Pg==",
                            "sender_device_uid": "%s",
                            "recipient_device_id": %d,
                            "userid": "%s",
                            "recipient_id": %d,
                            "environment_id": "that is not an environment id",
                            "update_environment_version": true
                        }
                    ]
                }
                `,
					users["admin"].Devices[0].UID,
					users["devops"].Devices[0].ID,
					users["admin"].UserID,
					users["devops"].ID,
				))),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "no such recipient",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "messages": [
                        {
                            "payload": "PGVuY3J5cHRlZF9jb250ZW50Pg==",
                            "sender_device_uid": "%s",
                            "recipient_device_id": %d,
                            "userid": "%s",
                            "recipient_id": 1290380294890348,
                            "environment_id": "%s",
                            "update_environment_version": true
                        }
                    ]
                }
                `,
					users["admin"].Devices[0].UID,
					users["devops"].Devices[0].ID,
					users["admin"].UserID,
					devEnvironment.EnvironmentID,
				))),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "recipient cannot write or write",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "messages": [
                        {
                            "payload": "PGVuY3J5cHRlZF9jb250ZW50Pg==",
                            "sender_device_uid": "%s",
                            "recipient_device_id": %d,
                            "userid": "%s",
                            "recipient_id": %d,
                            "environment_id": "%s",
                            "update_environment_version": true
                        }
                    ]
                }
                `,
					users["admin"].Devices[0].UID,
					users["developer"].Devices[0].ID,
					users["admin"].UserID,
					users["developer"].ID,
					prodEnvironment.EnvironmentID,
				))),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want: &models.GetEnvironmentsResponse{
				Environments: []models.Environment{},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "sender device not found",
			args: args{
				in0: router.Params{},
				body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "messages": [
                        {
                            "payload": "PGVuY3J5cHRlZF9jb250ZW50Pg==",
                            "sender_device_uid": "this is not a device",
                            "recipient_device_id": %d,
                            "userid": "%s",
                            "recipient_id": %d,
                            "environment_id": "%s",
                            "update_environment_version": true
                        }
                    ]
                }
                `,
					users["developer"].Devices[0].ID,
					users["admin"].UserID,
					users["developer"].ID,
					devEnvironment.EnvironmentID,
				))),
				Repo: newFakeRepo(),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "no device",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := WriteMessages(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"WriteMessages() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"WriteMessages() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			if tt.want == nil {
				if got.(*models.GetEnvironmentsResponse) != nil {
					t.Errorf("WriteMessages() got = %v, want nil", got)
				}
				return
			}

			gotResponse := got.(*models.GetEnvironmentsResponse)
			gotEnvironments := gotResponse.Environments
			wantEnvironments := tt.want.Environments
			gotLen := len(gotEnvironments)
			wantLen := len(wantEnvironments)
			if gotLen != wantLen {
				t.Errorf("WriteMessages() got = %v, want %v", got, tt.want)
				return
			}

			for index, wantEnvironment := range wantEnvironments {
				gotEnvironment := gotEnvironments[index]
				if gotEnvironment.EnvironmentID != wantEnvironment.EnvironmentID {
					t.Errorf("WriteMessages() got = %v, want %v", got, tt.want)
				}
				if gotEnvironment.VersionID == wantEnvironment.VersionID {
					t.Errorf(
						"WriteMessages() got VersionID %v want a different one",
						gotEnvironment.VersionID,
					)
				}
			}
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	project, users, messages := seedMessages(false)
	defer teardownMessages(project, users, messages)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *GenericResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "deletes a message",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"messageID": strconv.Itoa(int(messages[0].ID)),
				}),
				in1:  nil,
				Repo: newFakeRepo(),
				user: models.User{},
			},
			want: &GenericResponse{
				Success: true,
			},
			wantStatus: http.StatusNoContent,
			wantErr:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := DeleteMessage(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"DeleteMessage() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"DeleteMessage() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.(*GenericResponse)
			if !reflect.DeepEqual(gotResponse, tt.want) {
				t.Errorf("DeleteMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteExpiredMessages(t *testing.T) {
	messages := seedExpiredMessages()
	defer teardownExpiredMessages(messages)

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "deletes the messages",
			args: args{
				w: &fakeHttpResponse{
					header: map[string][]string{},
					body:   nil,
					status: 0,
				},
				r: &http.Request{
					Header: map[string][]string{
						"X-Ks-Ttl": {"a very long secret header"},
					},
				},
				in2: []httprouter.Param{},
			},
			wantStatus: http.StatusOK,
			wantBody:   "OK\n",
		},
		{
			name: "auth failure",
			args: args{
				w: &fakeHttpResponse{
					header: map[string][]string{},
					body:   nil,
					status: 0,
				},
				r: &http.Request{
					Header: map[string][]string{
						"X-Ks-Ttl": {"a very bad secret header"},
					},
				},
				in2: []httprouter.Param{},
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Unauthorized\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteExpiredMessages(tt.args.w, tt.args.r, tt.args.in2)
			resp := tt.args.w.(*fakeHttpResponse)
			body := resp.BodyString()
			if body != tt.wantBody {
				t.Errorf(
					"DeleteExpiredMessages() got body %v, want %v",
					body,
					tt.wantBody,
				)
			}
			if resp.status != tt.wantStatus {
				t.Errorf(
					"DeleteExpiredMessages() got status %v, want %v",
					resp.status,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestAlertMessagesWillExpire(t *testing.T) {
	messages := seedSoonToExpireMessages()
	defer teardownExpiredMessages(messages)

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "alerts messages will expire",
			args: args{
				w: &fakeHttpResponse{
					header: map[string][]string{},
					body:   nil,
					status: http.StatusOK,
				},
				r: &http.Request{
					Header: map[string][]string{
						"X-Ks-Ttl": {"a very long secret header"},
					},
				},
				in2: []httprouter.Param{},
			},
			wantStatus: http.StatusOK,
			wantBody:   "OK\n",
		},
		{
			name: "auth failure",
			args: args{
				w: &fakeHttpResponse{
					header: map[string][]string{},
					body:   nil,
					status: http.StatusOK,
				},
				r: &http.Request{
					Header: map[string][]string{
						"X-Ks-Ttl": {"a very bad secret header"},
					},
				},
				in2: []httprouter.Param{},
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Unauthorized\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AlertMessagesWillExpire(tt.args.w, tt.args.r, tt.args.in2)
			resp := tt.args.w.(*fakeHttpResponse)
			gotBody := resp.BodyString()

			if resp.status != tt.wantStatus {
				t.Errorf(
					"AlertMessagesWillExpire() got status = %v, want %v",
					resp.status,
					tt.wantStatus,
				)
			}

			if gotBody != tt.wantBody {
				t.Errorf(
					"AlertMessagesWillExpire() got = %v, want %v",
					gotBody,
					tt.wantBody,
				)
			}
		})
	}
}

func seedMessages(paid bool) (
	project models.Project,
	users map[string]models.User,
	messages []models.Message,
) {
	users = make(map[string]models.User)

	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		faker.FakeData(&project)
		db.Create(&project)

		orga := models.Organization{}
		faker.FakeData(&orga)
		if paid {
			orga.Paid = true
			db.Create(&orga)
		}
		project.Organization = orga
		project.OrganizationID = orga.ID
		db.Save(&project)

		roles := make([]models.Role, 0)
		db.Find(&roles)

		adminRole := models.Role{}
		db.Where("name = ?", "admin").First(&adminRole)

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
			db.Create(&env)

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
			db.Create(&user)

			device := models.Device{}
			faker.FakeData(&device)
			db.Create(&device)

			db.Model(&user).Association("Devices").Append(&device)
			user.Devices = []models.Device{device}

			roleID := role.ID
			if !orga.Paid {
				roleID = adminRole.ID
			}

			projectMember := models.ProjectMember{
				UserID:    user.ID,
				ProjectID: project.ID,
				RoleID:    roleID,
			}

			db.Create(&projectMember)
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

			db.Create(&message)

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
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(&messages)
		db.Delete(&project.Environments)

		if db.Error != nil {
			return db.Error
		}

		for _, user := range users {
			db.Model(&user).Association("Devices").Delete(&user.Devices)
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
	if err != nil {
		panic(err)
	}
}

func seedExpiredMessages() []models.Message {
	messages := []models.Message{}

	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		for i := 0; i < 10; i++ {
			message := models.Message{}
			faker.FakeData(&message)

			db.Create(&message)
			message.CreatedAt = time.Now().Add(-7 * 24 * time.Hour)
			db.Save(&message)

			messages = append(messages, message)
		}

		return db.Error
	})

	return messages
}

func teardownExpiredMessages(messages []models.Message) {
	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(messages)
		return db.Error
	})
}

func seedSoonToExpireMessages() []models.Message {
	messages := []models.Message{}

	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		for i := 0; i < 10; i++ {
			message := models.Message{}
			faker.FakeData(&message)

			db.Create(&message)
			message.CreatedAt = time.Now().Add(-6 * 24 * time.Hour)
			db.Save(&message)

			messages = append(messages, message)
		}

		return db.Error
	})

	return messages
}

// +------ Fake Repo

type fakeRepo struct {
	err error
	inner repo.IRepo
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		inner: repo.NewRepo(),
	}
}

func (f *fakeRepo) CreateEnvironment(
	environment *models.Environment,
) repo.IRepo {
	f.inner.CreateEnvironment(environment)
	return f
}

func (f *fakeRepo) CreateEnvironmentType(
	environmentType *models.EnvironmentType,
) repo.IRepo {
	f.inner.CreateEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) CreateLoginRequest() models.LoginRequest {
	return f.inner.CreateLoginRequest()
}

func (f *fakeRepo) CreateProjectMember(
	projectMember *models.ProjectMember,
	role *models.Role,
) repo.IRepo {
	f.inner.CreateProjectMember(projectMember, role)
	return f
}

func (f *fakeRepo) CreateRole(role *models.Role) repo.IRepo {
	f.inner.CreateRole(role)
	return f
}

func (f *fakeRepo) CreateRoleEnvironmentType(
	rolesEnvironmentType *models.RolesEnvironmentType,
) repo.IRepo {
	f.inner.CreateRoleEnvironmentType(rolesEnvironmentType)
	return f
}

func (f *fakeRepo) DeleteLoginRequest(id string) bool {
	return f.inner.DeleteLoginRequest(id)
}

func (f *fakeRepo) DeleteAllProjectMembers(project *models.Project) repo.IRepo {
	f.inner.DeleteAllProjectMembers(project)
	return f
}

func (f *fakeRepo) DeleteExpiredMessages() repo.IRepo {
	f.inner.DeleteExpiredMessages()
	return f
}

func (f *fakeRepo) GetGroupedMessagesWillExpireByUser(
	groupedMessageUser *map[uint]emailer.GroupedMessagesUser,
) repo.IRepo {
	f.inner.GetGroupedMessagesWillExpireByUser(groupedMessageUser)
	return f
}

func (f *fakeRepo) DeleteMessage(messageID uint, userID uint) repo.IRepo {
	f.inner.DeleteMessage(messageID, userID)
	return f
}

func (f *fakeRepo) DeleteProject(project *models.Project) repo.IRepo {
	f.inner.DeleteProject(project)
	return f
}

func (f *fakeRepo) DeleteProjectsEnvironments(
	project *models.Project,
) repo.IRepo {
	f.inner.DeleteProjectsEnvironments(project)
	return f
}

func (f *fakeRepo) Err() error {
	if f.err != nil {
		return f.err
	}

	return f.inner.Err()
}

func (f *fakeRepo) FindUsers(
	userIDs []string,
	users *map[string]models.User,
	notFounds *[]string,
) repo.IRepo {
	f.inner.FindUsers(userIDs, users, notFounds)
	return f
}

func (f *fakeRepo) GetActivityLogs(
	projectID string,
	options models.GetLogsOptions,
	logs *[]models.ActivityLog,
) repo.IRepo {
	f.inner.GetActivityLogs(projectID, options, logs)
	return f
}

func (f *fakeRepo) GetChildrenRoles(
	role models.Role,
	roles *[]models.Role,
) repo.IRepo {
	f.inner.GetChildrenRoles(role, roles)
	return f
}

func (f *fakeRepo) GetDb() *gorm.DB {
	return f.inner.GetDb()
}

func (f *fakeRepo) GetEnvironment(environment *models.Environment) repo.IRepo {
	f.inner.GetEnvironment(environment)
	return f
}

func (f *fakeRepo) GetEnvironmentPublicKeys(
	envID string,
	publicKeys *models.PublicKeys,
) repo.IRepo {
	f.inner.GetEnvironmentPublicKeys(envID, publicKeys)
	return f
}

func (f *fakeRepo) GetEnvironmentType(
	environmentType *models.EnvironmentType,
) repo.IRepo {
	f.inner.GetEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) GetEnvironmentsByProjectUUID(
	projectUUID string,
	foundEnvironments *[]models.Environment,
) repo.IRepo {
	f.inner.GetEnvironmentsByProjectUUID(projectUUID, foundEnvironments)
	return f
}

func (f *fakeRepo) GetInvitableRoles(
	role models.Role,
	roles *[]models.Role,
) repo.IRepo {
	f.inner.GetInvitableRoles(role, roles)
	return f
}

func (f *fakeRepo) GetLoginRequest(
	loginRequest string,
) (models.LoginRequest, bool) {
	return f.inner.GetLoginRequest(loginRequest)
}

func (f *fakeRepo) GetMessage(message *models.Message) repo.IRepo {
	f.inner.GetMessage(message)
	return f
}

func (f *fakeRepo) GetMessagesForUserOnEnvironment(
	device models.Device,
	environment models.Environment,
	message *models.Message,
) repo.IRepo {
	f.inner.GetMessagesForUserOnEnvironment(device, environment, message)
	return f
}

func (f *fakeRepo) GetOrCreateEnvironment(
	environment *models.Environment,
) repo.IRepo {
	f.inner.GetOrCreateEnvironment(environment)
	return f
}

func (f *fakeRepo) GetOrCreateEnvironmentType(
	environmentType *models.EnvironmentType,
) repo.IRepo {
	f.inner.GetOrCreateEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) GetOrCreateProject(project *models.Project) repo.IRepo {
	f.inner.GetOrCreateProject(project)
	return f
}

func (f *fakeRepo) GetOrCreateProjectMember(
	projectMember *models.ProjectMember,
	iRepo string,
) repo.IRepo {
	f.inner.GetOrCreateProjectMember(projectMember, iRepo)
	return f
}

func (f *fakeRepo) GetOrCreateRole(role *models.Role) repo.IRepo {
	f.inner.GetOrCreateRole(role)
	return f
}

func (f *fakeRepo) GetOrCreateRoleEnvType(
	rolesEnvironmentType *models.RolesEnvironmentType,
) repo.IRepo {
	f.inner.GetOrCreateRoleEnvType(rolesEnvironmentType)
	return f
}

func (f *fakeRepo) GetOrCreateUser(user *models.User) repo.IRepo {
	f.inner.GetOrCreateUser(user)
	return f
}

func (f *fakeRepo) GetProject(project *models.Project) repo.IRepo {
	if project.UUID == "crash it" {
			f.err = errors.New("unexpected error")
			return f
	}

	f.inner.GetProject(project)
	return f
}

func (f *fakeRepo) GetProjectByUUID(
	uuid string,
	project *models.Project,
) repo.IRepo {
	f.inner.GetProjectByUUID(uuid, project)
	return f
}

func (f *fakeRepo) GetProjectMember(
	projectMember *models.ProjectMember,
) repo.IRepo {
	f.inner.GetProjectMember(projectMember)
	return f
}

func (f *fakeRepo) GetProjectsOrganization(
	id string,
	organization *models.Organization,
) repo.IRepo {
	f.inner.GetProjectsOrganization(id, organization)
	return f
}

func (f *fakeRepo) OrganizationCountMembers(
	organization *models.Organization,
	iRepo *int64,
) repo.IRepo {
	f.inner.OrganizationCountMembers(organization, iRepo)
	return f
}

func (f *fakeRepo) GetRole(role *models.Role) repo.IRepo {
	f.inner.GetRole(role)
	return f
}

func (f *fakeRepo) GetRoles(role *[]models.Role) repo.IRepo {
	f.inner.GetRoles(role)
	return f
}

func (f *fakeRepo) GetRolesEnvironmentType(
	rolesEnvironmentType *models.RolesEnvironmentType,
) repo.IRepo {
	f.inner.GetRolesEnvironmentType(rolesEnvironmentType)
	return f
}

func (f *fakeRepo) GetRolesMemberCanInvite(
	projectMember models.ProjectMember,
	roles *[]models.Role,
) repo.IRepo {
	f.inner.GetRolesMemberCanInvite(projectMember, roles)
	return f
}

func (f *fakeRepo) GetUser(user *models.User) repo.IRepo {
	f.inner.GetUser(user)
	return f
}

func (f *fakeRepo) GetUserByEmail(id string, user *[]models.User) repo.IRepo {
	f.inner.GetUserByEmail(id, user)
	return f
}

func (f *fakeRepo) IsMemberOfProject(
	project *models.Project,
	projectMember *models.ProjectMember,
) repo.IRepo {
	f.inner.IsMemberOfProject(project, projectMember)
	return f
}

func (f *fakeRepo) ListProjectMembers(
	userIDList []string,
	projectMember *[]models.ProjectMember,
) repo.IRepo {
	f.inner.ListProjectMembers(userIDList, projectMember)
	return f
}

func (f *fakeRepo) MessageService() *message.MessageService {
	return f.inner.MessageService()
}

func (f *fakeRepo) ProjectAddMembers(
	project models.Project,
	memberRole []models.MemberRole,
	user models.User,
) repo.IRepo {
	f.inner.ProjectAddMembers(project, memberRole, user)
	return f
}

func (f *fakeRepo) UsersInMemberRoles(
	mers []models.MemberRole,
) (map[string]models.User, []string) {
	return f.inner.UsersInMemberRoles(mers)
}

func (f *fakeRepo) SetNewlyCreatedDevice(
	flag bool,
	deviceID uint,
	userID uint,
) repo.IRepo {
	f.inner.SetNewlyCreatedDevice(flag, deviceID, userID)
	return f
}

func (f *fakeRepo) ProjectGetAdmins(
	project *models.Project,
	members *[]models.ProjectMember,
) repo.IRepo {
	f.inner.ProjectGetAdmins(project, members)
	return f
}

func (f *fakeRepo) ProjectIsMemberAdmin(
	project *models.Project,
	member *models.ProjectMember,
) bool {
	return f.inner.ProjectIsMemberAdmin(project, member)
}

func (f *fakeRepo) ProjectGetMembers(
	project *models.Project,
	projectMember *[]models.ProjectMember,
) repo.IRepo {
	f.inner.ProjectGetMembers(project, projectMember)
	return f
}

func (f *fakeRepo) ProjectLoadUsers(project *models.Project) repo.IRepo {
	f.inner.ProjectLoadUsers(project)
	return f
}

func (f *fakeRepo) ProjectRemoveMembers(
	project models.Project,
	iRepo []string,
) repo.IRepo {
	f.inner.ProjectRemoveMembers(project, iRepo)
	return f
}

func (f *fakeRepo) ProjectSetRoleForUser(
	projet models.Project,
	user models.User,
	role models.Role,
) repo.IRepo {
	f.inner.ProjectSetRoleForUser(projet, user, role)
	return f
}

func (f *fakeRepo) CheckMembersAreInProject(
	project models.Project,
	members []string,
) ([]string, error) {
	return f.inner.CheckMembersAreInProject(project, members)
}

func (f *fakeRepo) RemoveOldMessageForRecipient(
	userID uint,
	environmentID string,
) repo.IRepo {
	f.inner.RemoveOldMessageForRecipient(userID, environmentID)
	return f
}

func (f *fakeRepo) SaveActivityLog(al *models.ActivityLog) repo.IRepo {
	f.inner.SaveActivityLog(al)
	return f
}

func (f *fakeRepo) SetLoginRequestCode(
	code string,
	c string,
) models.LoginRequest {
	return f.inner.SetLoginRequestCode(code, c)
}

func (f *fakeRepo) SetNewVersionID(environment *models.Environment) error {
	return f.inner.SetNewVersionID(environment)
}

func (f *fakeRepo) WriteMessage(
	user models.User,
	message models.Message,
) repo.IRepo {
	f.inner.WriteMessage(user, message)
	return f
}

func (f *fakeRepo) GetDevices(id uint, device *[]models.Device) repo.IRepo {
	f.inner.GetDevices(id, device)
	return f
}

func (f *fakeRepo) GetNewlyCreatedDevices(device *[]models.Device) repo.IRepo {
	f.inner.GetNewlyCreatedDevices(device)
	return f
}

func (f *fakeRepo) GetDevice(device *models.Device) repo.IRepo {
	f.inner.GetDevice(device)
	return f
}

func (f *fakeRepo) GetDeviceByUserID(
	userID uint,
	device *models.Device,
) repo.IRepo {
	f.inner.GetDeviceByUserID(userID, device)
	return f
}

func (f *fakeRepo) UpdateDeviceLastUsedAt(deviceUID string) repo.IRepo {
	f.inner.UpdateDeviceLastUsedAt(deviceUID)
	return f
}

func (f *fakeRepo) RevokeDevice(userID uint, deviceUID string) repo.IRepo {
	f.inner.RevokeDevice(userID, deviceUID)
	return f
}

func (f *fakeRepo) GetAdminsFromUserProjects(
	userID uint,
	adminProjectsMap *map[string][]string,
) repo.IRepo {
	f.inner.GetAdminsFromUserProjects(userID, adminProjectsMap)
	return f
}

func (f *fakeRepo) CreateOrganization(orga *models.Organization) repo.IRepo {
	f.inner.CreateOrganization(orga)
	return f
}

func (f *fakeRepo) UpdateOrganization(orga *models.Organization) repo.IRepo {
	f.inner.UpdateOrganization(orga)
	return f
}

func (f *fakeRepo) OrganizationSetCustomer(
	organization *models.Organization,
	customer string,
) repo.IRepo {
	f.inner.OrganizationSetCustomer(organization, customer)
	return f
}

func (f *fakeRepo) OrganizationSetSubscription(
	organization *models.Organization,
	subscription string,
) repo.IRepo {
	f.inner.OrganizationSetSubscription(organization, subscription)
	return f
}

func (f *fakeRepo) GetOrganization(orga *models.Organization) repo.IRepo {
	f.inner.GetOrganization(orga)
	return f
}

func (f *fakeRepo) GetOrganizations(
	userID uint,
	result *[]models.Organization,
) repo.IRepo {
	f.inner.GetOrganizations(userID, result)
	return f
}

func (f *fakeRepo) GetOwnedOrganizations(
	userID uint,
	result *[]models.Organization,
) repo.IRepo {
	f.inner.GetOwnedOrganizations(userID, result)
	return f
}

func (f *fakeRepo) GetOwnedOrganizationByName(
	userID uint,
	name string,
	orgas *[]models.Organization,
) repo.IRepo {
	f.inner.GetOwnedOrganizationByName(userID, name, orgas)
	return f
}

func (f *fakeRepo) GetOrganizationByName(
	userID uint,
	name string,
	orga *[]models.Organization,
) repo.IRepo {
	f.inner.GetOrganizationByName(userID, name, orga)
	return f
}

func (f *fakeRepo) GetOrganizationProjects(
	organization *models.Organization,
	project *[]models.Project,
) repo.IRepo {
	f.inner.GetOrganizationProjects(organization, project)
	return f
}

func (f *fakeRepo) GetOrganizationMembers(
	orgaID uint,
	result *[]models.ProjectMember,
) repo.IRepo {
	f.inner.GetOrganizationMembers(orgaID, result)
	return f
}

func (f *fakeRepo) IsUserOwnerOfOrga(
	user *models.User,
	organization *models.Organization,
) (bool, error) {
	return f.inner.IsUserOwnerOfOrga(user, organization)
}

func (f *fakeRepo) IsProjectOrganizationPaid(str string) (bool, error) {
	return f.inner.IsProjectOrganizationPaid(str)
}

func (f *fakeRepo) CreateCheckoutSession(
	checkoutSession *models.CheckoutSession,
) repo.IRepo {
	f.inner.CreateCheckoutSession(checkoutSession)
	return f
}

func (f *fakeRepo) GetCheckoutSession(
	str string,
	checkoutSession *models.CheckoutSession,
) repo.IRepo {
	f.inner.GetCheckoutSession(str, checkoutSession)
	return f
}

func (f *fakeRepo) UpdateCheckoutSession(
	checkoutSession *models.CheckoutSession,
) repo.IRepo {
	f.inner.UpdateCheckoutSession(checkoutSession)
	return f
}

func (f *fakeRepo) DeleteCheckoutSession(
	checkoutSession *models.CheckoutSession,
) repo.IRepo {
	f.inner.DeleteCheckoutSession(checkoutSession)
	return f
}

func (f *fakeRepo) OrganizationSetPaid(
	organization *models.Organization,
	paid bool,
) repo.IRepo {
	f.inner.OrganizationSetPaid(organization, paid)
	return f
}

func (f *fakeRepo) GetUserProjects(
	userID uint,
	projects *[]models.Project,
) repo.IRepo {
	f.inner.GetUserProjects(userID, projects)
	return f
}
