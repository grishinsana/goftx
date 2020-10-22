# goftx
FTX exchange golang library

### No Logged In Error
"Not logged in" errors usually come from a wrong signatures. FTX released an article on how to authenticate https://blog.ftx.com/blog/api-authentication/

If you have unauthorized error to private methods, then you need to use SetServerTimeDiff()
```go
ftx := New()
ftx.SetServerTimeDiff()
```
