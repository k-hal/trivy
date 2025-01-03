package poetry

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
)

func Test_poetryLibraryAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want *analyzer.AnalysisResult
	}{
		{
			name: "happy path",
			dir:  "testdata/happy",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Poetry,
						FilePath: "poetry.lock",
						Packages: types.Packages{
							{
								ID:           "certifi@2022.12.7",
								Name:         "certifi",
								Version:      "2022.12.7",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "charset-normalizer@2.1.1",
								Name:         "charset-normalizer",
								Version:      "2.1.1",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "click@7.1.2",
								Name:         "click",
								Version:      "7.1.2",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "flask@1.1.4",
								Name:         "flask",
								Version:      "1.1.4",
								Relationship: types.RelationshipDirect,
								DependsOn: []string{
									"click@7.1.2",
									"itsdangerous@1.1.0",
									"jinja2@2.11.3",
									"werkzeug@1.0.1",
								},
							},
							{
								ID:           "idna@3.4",
								Name:         "idna",
								Version:      "3.4",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "itsdangerous@1.1.0",
								Name:         "itsdangerous",
								Version:      "1.1.0",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "jinja2@2.11.3",
								Name:         "jinja2",
								Version:      "2.11.3",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
								DependsOn: []string{
									"markupsafe@2.1.2",
								},
							},
							{
								ID:           "markupsafe@2.1.2",
								Name:         "markupsafe",
								Version:      "2.1.2",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "requests@2.28.1",
								Name:         "requests",
								Version:      "2.28.1",
								Relationship: types.RelationshipDirect,
								DependsOn: []string{
									"certifi@2022.12.7",
									"charset-normalizer@2.1.1",
									"idna@3.4",
									"urllib3@1.26.14",
								},
							},
							{
								ID:           "urllib3@1.26.14",
								Name:         "urllib3",
								Version:      "1.26.14",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "werkzeug@1.0.1",
								Name:         "werkzeug",
								Version:      "1.0.1",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
						},
					},
				},
			},
		},
		{
			name: "no pyproject.toml",
			dir:  "testdata/no-pyproject",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Poetry,
						FilePath: "poetry.lock",
						Packages: types.Packages{
							{
								ID:      "click@8.1.3",
								Name:    "click",
								Version: "8.1.3",
								DependsOn: []string{
									"colorama@0.4.6",
								},
							},
							{
								ID:      "colorama@0.4.6",
								Name:    "colorama",
								Version: "0.4.6",
							},
						},
					},
				},
			},
		},
		{
			name: "wrong pyproject.toml",
			dir:  "testdata/wrong-pyproject",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Poetry,
						FilePath: "poetry.lock",
						Packages: types.Packages{
							{
								ID:      "click@8.1.3",
								Name:    "click",
								Version: "8.1.3",
								DependsOn: []string{
									"colorama@0.4.6",
								},
							},
							{
								ID:      "colorama@0.4.6",
								Name:    "colorama",
								Version: "0.4.6",
							},
						},
					},
				},
			},
		},
		{
			name: "broken poetry.lock",
			dir:  "testdata/sad",
			want: &analyzer.AnalysisResult{},
		},
		{
			// docker run --name poetry --rm -it python@sha256:e1141f10176d74d1a0e87a7c0a0a5a98dd98ec5ac12ce867768f40c6feae2fd9 sh
			// wget -qO- https://install.python-poetry.org | POETRY_VERSION=1.8.5 python3 -
			// export PATH="/root/.local/bin:$PATH"
			// poetry new groups && cd groups
			// poetry add requests@2.32.3
			// poetry add --group dev pytest@8.3.4
			// poetry add --group lint ruff@0.8.3
			// poetry add --optional typing-inspect@0.9.0
			name: "skip deps from groups",
			dir:  "testdata/with-groups",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Poetry,
						FilePath: "poetry.lock",
						Packages: types.Packages{
							{
								ID:           "certifi@2024.12.14",
								Name:         "certifi",
								Version:      "2024.12.14",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "charset-normalizer@3.4.0",
								Name:         "charset-normalizer",
								Version:      "3.4.0",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "idna@3.10",
								Name:         "idna",
								Version:      "3.10",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "mypy-extensions@1.0.0",
								Name:         "mypy-extensions",
								Version:      "1.0.0",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:      "requests@2.32.3",
								Name:    "requests",
								Version: "2.32.3",
								DependsOn: []string{
									"certifi@2024.12.14",
									"charset-normalizer@3.4.0",
									"idna@3.10",
									"urllib3@2.2.3",
								},
								Relationship: types.RelationshipDirect,
							},
							{
								ID:           "ruff@0.8.3",
								Name:         "ruff",
								Version:      "0.8.3",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:           "typing-extensions@4.12.2",
								Name:         "typing-extensions",
								Version:      "4.12.2",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
							{
								ID:      "typing-inspect@0.9.0",
								Name:    "typing-inspect",
								Version: "0.9.0",
								DependsOn: []string{
									"mypy-extensions@1.0.0",
									"typing-extensions@4.12.2",
								},
								Relationship: types.RelationshipDirect,
							},
							{
								ID:           "urllib3@2.2.3",
								Name:         "urllib3",
								Version:      "2.2.3",
								Indirect:     true,
								Relationship: types.RelationshipIndirect,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := newPoetryAnalyzer(analyzer.AnalyzerOptions{})
			require.NoError(t, err)

			got, err := a.PostAnalyze(context.Background(), analyzer.PostAnalysisInput{
				FS: os.DirFS(tt.dir),
			})

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
