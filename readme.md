###################### Command Run App ######################

1. Command Run App => CompileDaemon -build="go build -o app.exe" -command="./app.exe"
2. Create file .env for keep secret key.
    GOOGLE_CLIENT_ID=your client id
    GOOGLE_CLIENT_SECRET=your secret key
    GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
    GOOGLE_USERINFO=https://www.googleapis.com/oauth2/v2/userinfo
    GOOGLE_USERINFO_EMAIL=https://www.googleapis.com/auth/userinfo.email
    GOOGLE_USERINFO_PROFILE=https://www.googleapis.com/auth/userinfo.profile
3. Create file rsa.pem and in file and private key for system middleware.