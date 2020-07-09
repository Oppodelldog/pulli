package git

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGetStats(t *testing.T) {
	testCases := []struct {
		gitOutput string
		want      PullStats
	}{
		{
			gitOutput: " 1 file changed, 1 insertion(+)",
			want: PullStats{
				Files:      1,
				Insertions: 1,
				Deletions:  0},
		},
		{
			gitOutput: " 2 files changed, 2 insertions(+)",
			want: PullStats{
				Files:      2,
				Insertions: 2,
				Deletions:  0},
		},
		{
			gitOutput: " 1 file changed, 1 deletion(-)",
			want: PullStats{
				Files:      1,
				Insertions: 0,
				Deletions:  1},
		},
		{
			gitOutput: " 1 file changed, 2 deletions(-)",
			want: PullStats{
				Files:      1,
				Insertions: 0,
				Deletions:  2},
		},
		{
			gitOutput: " 2 files changed, 1 insertion(+), 1 deletion(-)`",
			want: PullStats{
				Files:      2,
				Insertions: 1,
				Deletions:  1},
		},
		{
			gitOutput: " 2 files changed, 2 insertions(+), 2 deletions(-)`",
			want: PullStats{
				Files:      2,
				Insertions: 2,
				Deletions:  2},
		},
	}

	for i, testCase := range testCases {
		t.Run("test index "+strconv.Itoa(i), func(t *testing.T) {
			got := GetStats(testCase.gitOutput)
			if !reflect.DeepEqual(testCase.want, got) {
				t.Fatalf("\nwant:\n%v\ngot:\n%v\n", testCase.want, got)
			}
		})
	}
}
