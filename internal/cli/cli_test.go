package cli

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedErr  bool
		expectedFlag Flags
		expectedArg  Args
	}{
		{
			name:        "基本参数",
			args:        []string{"https://github.com/owner/repo/issues/123"},
			expectedErr: false,
			expectedFlag: Flags{
				EnableReactions:   false,
				EnableUserLinks:   false,
			},
			expectedArg: Args{
				URL:        "https://github.com/owner/repo/issues/123",
				OutputFile: "",
			},
		},
		{
			name:        "带输出文件",
			args:        []string{"https://github.com/owner/repo/issues/123", "output.md"},
			expectedErr: false,
			expectedFlag: Flags{
				EnableReactions:   false,
				EnableUserLinks:   false,
			},
			expectedArg: Args{
				URL:        "https://github.com/owner/repo/issues/123",
				OutputFile: "output.md",
			},
		},
		{
			name:        "启用 Reactions",
			args:        []string{"-enable-reactions", "https://github.com/owner/repo/issues/123"},
			expectedErr: false,
			expectedFlag: Flags{
				EnableReactions:   true,
				EnableUserLinks:   false,
			},
			expectedArg: Args{
				URL:        "https://github.com/owner/repo/issues/123",
				OutputFile: "",
			},
		},
		{
			name:        "启用用户链接",
			args:        []string{"-enable-user-links", "https://github.com/owner/repo/issues/123"},
			expectedErr: false,
			expectedFlag: Flags{
				EnableReactions:   false,
				EnableUserLinks:   true,
			},
			expectedArg: Args{
				URL:        "https://github.com/owner/repo/issues/123",
				OutputFile: "",
			},
		},
		{
			name:        "所有标志",
			args:        []string{"-enable-reactions", "-enable-user-links", "https://github.com/owner/repo/issues/123", "output.md"},
			expectedErr: false,
			expectedFlag: Flags{
				EnableReactions:   true,
				EnableUserLinks:   true,
			},
			expectedArg: Args{
				URL:        "https://github.com/owner/repo/issues/123",
				OutputFile: "output.md",
			},
		},
		{
			name:        "缺少参数",
			args:        []string{},
			expectedErr: true,
		},
		{
			name:        "未知标志",
			args:        []string{"-unknown", "https://github.com/owner/repo/issues/123"},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, args, err := ParseArgs(tt.args)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("ParseArgs(%v) expected error, got nil", tt.args)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseArgs(%v) unexpected error: %v", tt.args, err)
				return
			}

			if flags.EnableReactions != tt.expectedFlag.EnableReactions {
				t.Errorf("ParseArgs(%v).EnableReactions = %v, want %v", tt.args, flags.EnableReactions, tt.expectedFlag.EnableReactions)
			}
			if flags.EnableUserLinks != tt.expectedFlag.EnableUserLinks {
				t.Errorf("ParseArgs(%v).EnableUserLinks = %v, want %v", tt.args, flags.EnableUserLinks, tt.expectedFlag.EnableUserLinks)
			}
			if args.URL != tt.expectedArg.URL {
				t.Errorf("ParseArgs(%v).URL = %q, want %q", tt.args, args.URL, tt.expectedArg.URL)
			}
			if args.OutputFile != tt.expectedArg.OutputFile {
				t.Errorf("ParseArgs(%v).OutputFile = %q, want %q", tt.args, args.OutputFile, tt.expectedArg.OutputFile)
			}
		})
	}
}
