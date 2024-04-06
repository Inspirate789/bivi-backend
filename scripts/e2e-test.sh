set -x

if [ "$1" = "docker" ]; then
  reports_path=/test-reports
else
  reports_path=$(pwd)/test-reports
fi

go test -c -v ./internal/pkg/app/e2e_test.go
export ALLURE_OUTPUT_PATH="$reports_path" && ./app.test
