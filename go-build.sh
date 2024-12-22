#!/bin/bash

# 설정
APP_NAME="libraryBackend"
BUILD_DIR="./build"
REMOTE_USER="ubuntu"
REMOTE_HOST="152.67.208.135"
REMOTE_DIR="."
SSH_KEY="~/.key/oracle-24-12-22.key"

REMOTE_COMMAND='
process_info=$(sudo lsof -i :80 | grep libraryBa);
if [ -z "$process_info" ]; then
  echo "libraryBackend 관련 프로세스가 포트 80에서 실행 중이지 않습니다.";
  exit 0;
fi;
pid=$(echo "$process_info" | awk "{print \$2}");
if [ -n "$pid" ]; then
  echo "PID $pid를 종료합니다.";
  sudo kill -9 $pid;
  echo "프로세스가 종료되었습니다.";
else
  echo "PID를 찾을 수 없습니다. 종료 작업을 수행하지 않았습니다.";
fi
'


# # 1. 바이너리 빌드
# echo "🔨 빌드 중..."
# GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/$APP_NAME"
# if [ $? -ne 0 ]; then
#     echo "❌ 빌드 실패"
#     exit 1
# fi
# echo "✅ 빌드 완료: $BUILD_DIR/$APP_NAME"

ssh -i $SSH_KEY "$REMOTE_USER@$REMOTE_HOST" $REMOTE_COMMAND

# 2. 원격 서버로 복사
echo "📦 원격 서버로 파일 복사 중..."
scp -i $SSH_KEY "$BUILD_DIR/$APP_NAME" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR"
if [ $? -ne 0 ]; then
    echo "❌ 파일 복사 실패"
    exit 1
fi
echo "✅ 파일 복사 완료"

# 3. 원격 서버에서 실행 권한 부여 및 백그라운드 실행
echo "🔧 실행 권한 부여 및 백그라운드 실행 중..."
# -t: TTY(가상 터미널)를 강제로 할당해 명령 실행이 끝난 후 SSH 세션을 닫음.
ssh -i $SSH_KEY "$REMOTE_USER@$REMOTE_HOST" "chmod +x ./libraryBackend"
ssh -i $SSH_KEY "$REMOTE_USER@$REMOTE_HOST" "nohup sudo $REMOTE_DIR/$APP_NAME > $REMOTE_DIR/$APP_NAME.log 2>&1 &"

echo "✅ 실행 완료: $REMOTE_DIR/$APP_NAME (로그: $REMOTE_DIR/$APP_NAME.log)"

echo "🚀 배포 완료!"