# GDH 인증 서버 입니다.

## 환경변수

> 다음 2가지 환경변수를 설정해야 합니다.

```dotenv
DB_URL="postgresql://ID:PW@호스트주소:포트번호/데이터베이스명"
JWT_SECRET="JWT_SECRET_STRING"
```

## 포트

- 50501 (grpc)