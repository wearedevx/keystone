package message

import (
	"reflect"
	"testing"
)

func TestNewMessageService(t *testing.T) {
	t.Run("instanciate a message service", func(t *testing.T) {
		service := NewMessageService().(*messageService)

		if service == nil {
			t.Errorf("NewMessageService(): got nil")
		}

		if service.redis == nil {
			t.Errorf("NewMessageService(): service redis is nil")
		}
	})
}

func TestMessageService_GetMessageByUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		fixtures map[string]string
		name     string
		args     args
		want     []byte
		wantErr  bool
	}{
		{
			name: "finds the message",
			fixtures: map[string]string{
				"uuid": "value",
			},
			args: args{
				uuid: "uuid",
			},
			want:    []byte("value"),
			wantErr: false,
		},
		{
			name: "does not find the message",
			fixtures: map[string]string{
				"uuid": "value",
			},
			args: args{
				uuid: "non existing value",
			},
			want:    []byte(""),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessageService().(*messageService)
			m.redis.SetupFixtures(tt.fixtures)
			got, err := m.GetMessageByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"MessageService.GetMessageByUUID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"MessageService.GetMessageByUUID() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestMessageService_WriteMessageWithUUID(t *testing.T) {
	type args struct {
		uuid  string
		value []byte
	}
	tests := []struct {
		name     string
		fixtures map[string]string
		args     args
		wantErr  bool
	}{
		{
			name:     "writes a message",
			fixtures: make(map[string]string),
			args: args{
				uuid:  "written-uuid",
				value: []byte("written-value"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessageService().(*messageService)
			m.redis.SetupFixtures(tt.fixtures)
			if err := m.WriteMessageWithUUID(tt.args.uuid, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf(
					"MessageService.WriteMessageWithUUID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestMessageService_DeleteMessageWithUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name     string
		fixtures map[string]string
		args     args
		wantErr  bool
	}{
		{
			name: "deletes a message",
			fixtures: map[string]string{
				"delete-uuid": "delete-value",
			},
			args: args{
				uuid: "delete-uuid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessageService().(*messageService)
			m.redis.SetupFixtures(tt.fixtures)
			if err := m.DeleteMessageWithUUID(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf(
					"MessageService.DeleteMessageWithUUID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if _, ok := tt.fixtures[tt.args.uuid]; ok {
				t.Errorf(
					"MessageService.DeleteMessageWithUUID() still has uuid %v",
					tt.args.uuid,
				)
			}
		})
	}
}
