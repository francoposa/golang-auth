Browser Login Flow with Unprotected Login Node Served From Backend

1. Browser click or redirect from a protected page creates request to authentication server /api/login/browser
2. Authentication server creates and persist login: id, timestamps, csrf token
3. Authentication server responds with redirect to Login UI with login id appended and csrf token set in cookie
4. Browser follows redirect to Login UI on same domain.
Cookies from authentication server's response are available on Login UI request/response, as it is the same domain
5. Login UI backend handler uses login id to fetch csrf token, as headers & body are not sent on redirect
6. Login UI backend handler responds with login form & csrf token in hidden form body element
7. 
