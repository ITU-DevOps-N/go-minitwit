<!doctype html>
<title> {{.title}} | MiniTwit</title>
<link rel=stylesheet type=text/css href="/web/static/style.css">

<div class=page>
  <h1>MiniTwit</h1>
  <div class=navigation>
  {{ if .user }}
    <a href="/timeline">My Timeline</a> |
    <a href="/public_timeline">Public Timeline</a> |
    <a href="/logout">sign out [{{ .user.username }}]</a>
  {{ else }}
    <a href="/public_timeline">public timeline</a> |
    <a href="/register">sign up</a> |
    <a href="/login">sign in</a>
  {{ end }}
  </div>
  <div class=body>
    <h2>Sign Up</h2>
    {{ if .error }}<div class=error><strong>Error:</strong> {{ .error }}</div>{{ end }}
    <form action="" method=post>
      <dl>
        <dt>Username:
        <dd><input type=text name=username size=30 value="{{ .form.username }}">
        <dt>E-Mail:
        <dd><input type=text name=email size=30 value="{{ .form.email }}">
        <dt>Password:
        <dd><input type=password name=password1 size=30>
        <dt>Password <small>(repeat)</small>:
        <dd><input type=password name=password2 size=30>
      </dl>
      <div class=actions><input type=submit value="Sign Up"></div>
    </form>
  </div>