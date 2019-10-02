package cmd

import (
	"fmt"
	"strings"
	"testing"
)

type testRepo struct {
	head    string
	headErr error

	spinOffErr      error
	capturedSpinOff string

	commitErr         error
	capturedCommitMsg string
}

func (r *testRepo) GetCurrentBranch() (string, error) {
	return r.head, r.headErr
}

func (r *testRepo) CheckoutSpinOffBranch(branch string) error {
	r.capturedSpinOff = branch
	return r.spinOffErr
}

func (r *testRepo) CommitAllChanged(msg string) error {
	r.capturedCommitMsg = msg
	return r.commitErr
}

func TestCommitAllChanges(t *testing.T) {
	tests := []struct {
		name          string
		repo          *testRepo
		wantSpinOff   string
		wantCommitMsg string
		wantErr       bool
	}{
		{
			name: "Switch to autosave branch and commit changes",
			repo: &testRepo{
				head:       "master",
				headErr:    nil,
				spinOffErr: nil,
				commitErr:  nil,
			},
			wantSpinOff:   "autosave/master",
			wantCommitMsg: "git-autosave at ",
			wantErr:       false,
		},
		{
			name: "Don't switch if already on autosave branch and commit changes",
			repo: &testRepo{
				head:       "autosave/master",
				headErr:    nil,
				spinOffErr: nil,
				commitErr:  nil,
			},
			wantSpinOff:   "",
			wantCommitMsg: "git-autosave at ",
			wantErr:       false,
		},
		{
			name: "Fail if can't check head name",
			repo: &testRepo{
				head:       "",
				headErr:    fmt.Errorf("can't check head"),
				spinOffErr: nil,
				commitErr:  nil,
			},
			wantSpinOff:   "",
			wantCommitMsg: "",
			wantErr:       true,
		},
		{
			name: "Fail if can't switch to spin off",
			repo: &testRepo{
				head:       "master",
				headErr:    nil,
				spinOffErr: fmt.Errorf("can't switch to spin off"),
				commitErr:  nil,
			},
			wantSpinOff:   "autosave/master",
			wantCommitMsg: "",
			wantErr:       true,
		},
		{
			name: "Fail if can't commit",
			repo: &testRepo{
				head:       "autosave/master",
				headErr:    nil,
				spinOffErr: nil,
				commitErr:  fmt.Errorf("can't commit"),
			},
			wantSpinOff:   "",
			wantCommitMsg: "git-autosave at ",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CommitAllChanges(tt.repo); (err != nil) != tt.wantErr {
				t.Errorf("CommitAllChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.repo.capturedSpinOff != tt.wantSpinOff {
				t.Errorf("CommitAllChanges() created a spinoff branch = '%s', wanted '%s'", tt.repo.capturedSpinOff, tt.wantSpinOff)
			}
			if !strings.HasPrefix(tt.repo.capturedCommitMsg, tt.wantCommitMsg) {
				t.Errorf("CommitAllChanges() created a commit = '%s', wanted it to start with '%s'", tt.repo.capturedCommitMsg, tt.wantCommitMsg)
			}
		})
	}
}
