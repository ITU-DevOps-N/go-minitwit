<!doctype html>
<title> {{.title}} | MiniTwit</title>
<link rel=stylesheet type=text/css href="/web/static/style.css">

<div class=page>
  <h1>MiniTwit</h1>
  <div class=navigation>
  {{ if .user }}
    <a href="/user_timeline">my timeline</a> |
    <a href="/public_timeline">public timeline</a> |
    <a href="/logout">logout</a>
  {{ else }}
    <a href="/public_timeline">public timeline</a> |
    <a href="/register">sign up</a> |
    <a href="/login">login</a>
  {{ end }}
  </div>
  <div class=body>
    <h2>{{ .title }}</h2>
  {{ if .private }}
    {{ if .user_timeline }}
      <div class=followstatus>
      {{ if .user_page }}
        This is you!
      {{ else if .followed }}
        You are currently following this user.
        <a class=unfollow href="/unfollow?username={{.user}}">Unfollow user</a>.
      {{ else }}
        You are not yet following this user.
        <a class=follow href="/follow?username={{.user}}">Follow user</a>.
      {{ end }}
      </div>
    {{ else}}
      <div class=twitbox>
        <h3>What's on your mind {{ .user }}?</h3>
        <form action="/add_message" method=post>
          <p><input type=text name=message size=60>
          <input type=submit value="Share">
        </form>
      </div>
    {{ end }}
  {{ end }}
  <ul class=messages>
  {{ range .messages }}
    <li><img src="http://www.gravatar.com/avatar/{{ .author | getUserId }}?d=identicon&amp;s=48"><p>
      <strong><a href="/user_timeline?username={{.author}}">{{ .author }}</a></strong>
      {{ .text }}
      <small>&mdash; {{ .created_at  | formatAsDate}}</small>
  {{ else }}
    <li><em>There's no message so far.</em>
  {{ end }}
  </ul>
  </div>
  <div class=footer>
    MiniTwit &mdash; A Gin Application
  </div>
</div>
