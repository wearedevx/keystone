//go:build test
// +build test

package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

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
		wantCalled       []string
	}{
		{
			name: "gets message for a user - dev env for dev user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
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
			name: "fails while fetching project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"GetProject": errors.New("unexpected error"),
				}),
				user: users["developer"],
			},
			want:             &models.GetMessageByEnvironmentResponse{},
			wantEnvironments: []string{},
			wantStatus:       http.StatusBadRequest,
			wantErr:          "bad request: unexpected error",
		},
		{
			name: "fails while getting permissions",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"GetProjectMember": errors.New("not found"),
				}),
				user: users["developer"],
			},
			want:             &models.GetMessageByEnvironmentResponse{},
			wantEnvironments: []string{},
			wantStatus:       http.StatusNotFound,
			wantErr:          "failed to get permission: not found",
		},
		{
			name: "fails while getting messages from db",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"GetMessagesForUserOnEnvironment": errors.New(
						"unexpected error",
					),
				}),
				user: users["developer"],
			},
			want:             &models.GetMessageByEnvironmentResponse{},
			wantEnvironments: []string{},
			wantStatus:       http.StatusBadRequest,
			wantErr:          "failed to get: unexpected error",
		},
		{
			name: "fails while deleting message after a redis error",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"MessageService.GetMessageByUUID": errors.New(
						"unexpected error",
					),
					"DeleteMessage": errors.New("not found"),
				}),
				user: users["developer"],
			},
			want:             &models.GetMessageByEnvironmentResponse{},
			wantEnvironments: []string{},
			wantStatus:       http.StatusOK,
			wantErr:          "",
			wantCalled:       []string{"DeleteMessage"},
		},
		{
			name: "deletes db messages if read form redis fails",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["developer"].Devices[0].UID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"MessageService.GetMessageByUUID": errors.New(
						"unexpected error",
					),
				}),
				user: users["developer"],
			},
			want:             &models.GetMessageByEnvironmentResponse{},
			wantEnvironments: []string{},
			wantStatus:       http.StatusOK,
			wantErr:          "",
			wantCalled:       []string{"DeleteMessage"},
		},
		{
			name: "gets message for a user - dev env for lead-dev user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
					"device":    users["lead-dev"].Devices[0].UID,
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
			tt.args.Repo.ClearErr()
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

			if tt.wantCalled != nil {
				gotCalled := tt.args.Repo.(*fakeRepo).called

				found := false
				for _, w := range tt.wantCalled {
					found = false

					for _, g := range gotCalled {
						if g == w {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("GetMessagesFromProjectByUser() got called %v, want %v", gotCalled, tt.wantCalled)
					}
				}
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "recipient cannot read or write",
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "no device",
		},
		{
			name: "error while checking free org has no non admin members",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganizationMembers": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       &models.GetEnvironmentsResponse{},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "organization without an admin: unexpected error",
		},
		{
			name: "error getting permissions",
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
				Repo: newFakeRepo(map[string]error{
					"GetRolesEnvironmentType": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       &models.GetEnvironmentsResponse{},
			wantStatus: http.StatusNotFound,
			wantErr:    "failed to get permission",
		},
		{
			name: "error removing previous messages",
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
				Repo: newFakeRepo(map[string]error{
					"RemoveOldMessageForRecipient": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       &models.GetEnvironmentsResponse{},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to delete: unexpected error",
		},
		{
			name: "fails when sender device is not found",
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
				Repo: newFakeRepo(map[string]error{
					"GetDevice": repo.ErrorNotFound,
				}),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "no device",
		},
		{
			name: "fails trying to get sender device",
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
				Repo: newFakeRepo(map[string]error{
					"GetDevice": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    "no device",
		},
		{
			name: "fails writing messages in db",
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
				Repo: newFakeRepo(map[string]error{
					"WriteMessage": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to write message: unexpected error",
		},
		{
			name: "fails setting version id because not found",
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
				Repo: newFakeRepo(map[string]error{
					"SetNewVersionID": repo.ErrorNotFound,
				}),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "fails setting version id because unexpected",
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
				Repo: newFakeRepo(map[string]error{
					"SetNewVersionID": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to set environment version: unexpected error",
		},
		{
			name: "fails writing message on redis",
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
				Repo: newFakeRepo(map[string]error{
					"MessageService.WriteMessageWithUUID": errors.New("unexpected error"),
				}),
				user: users["admin"],
			},
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to write message: unexpected error",
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
			name: "fails if id not an uint",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"messageID": "not a valid id",
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
				user: models.User{},
			},
			want: &GenericResponse{
				Success: false,
				Error:   "strconv.ParseUint: parsing \"not a valid id\": invalid syntax",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    "strconv.ParseUint: parsing \"not a valid id\": invalid syntax: invalid syntax",
		},
		{
			name: "message not found",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"messageID": strconv.Itoa(int(messages[0].ID)),
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"GetMessage": repo.ErrorNotFound,
				}),
				user: models.User{},
			},
			want: &GenericResponse{
				Success: true,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "delete fails in db",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"messageID": strconv.Itoa(int(messages[0].ID)),
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"DeleteMessage": errors.New("unexpected error"),
				}),
				user: models.User{},
			},
			want: &GenericResponse{
				Error: "unexpected error",
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to delete: unexpected error",
		},
		{
			name: "deletes a message",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"messageID": strconv.Itoa(int(messages[1].ID)),
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
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
			if tt.want.Success != gotResponse.Success {
				t.Errorf("DeleteMessage() got = %v, want %v", got, tt.want)
			}

			if tt.want.Error != gotResponse.Error {
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

	new(repo.Repo).GetDB().Transaction(func(db *gorm.DB) error {
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
	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
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

	new(repo.Repo).GetDB().Transaction(func(db *gorm.DB) error {
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
	new(repo.Repo).GetDB().Transaction(func(db *gorm.DB) error {
		db.Delete(messages)
		return db.Error
	})
}

func seedSoonToExpireMessages() []models.Message {
	messages := []models.Message{}

	new(repo.Repo).GetDB().Transaction(func(db *gorm.DB) error {
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
