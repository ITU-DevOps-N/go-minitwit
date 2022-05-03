<!doctype html>
<title> {{.title}} | MiniTwit</title>
<link rel=stylesheet type=text/css href="/web/static/style.css">
<div class=page>
  <h1>MiniTwit</h1>
  <div class=navigation>
  {{ if .user }}
    <a href="/timeline">my timeline</a> |
    <a href="/public_timeline">public timeline</a> |
    <a href="/logout">sign out [{{ .user.username }}]</a>
  {{ else }}
    <a href="/public_timeline">public timeline</a> |
    <a href="/register">sign up</a> |
    <a href="/login">sign in</a>
  {{ end }}
</div>
<div class=body>
  <h2>Sign In</h2>
  {{ if .error }}<div class=error><strong>Error:</strong> {{ .error }}</div>{{ end }}
  {{ if .ErrorTitle}}
    <div class=error><strong>{{.ErrorTitle}}</strong>{{.ErrorMessage}}</div>
  {{end}}
  <form action="" method=post>
    <dl>
      <dt>Username:
      <dd><input type=text name=username size=30 value="{{ .form.username }}">
      <dt>Password:
      <dd><input type=password name=password size=30>
    </dl>
    <div class=actions><input type=submit value="Sign In"></div>
  </form>
</div>