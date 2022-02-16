<!doctype html>
<title> {{.title}}| MiniTwit</title>
<link rel=stylesheet type=text/css href="/static/style.css">

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
    <h2>{{ .title }}</h2>
  {{ if .user }}
    {{ if .endpoint eq "user_timeline" }}
      <div class=followstatus>
      {{ if .user.user_id eq .profile_user.user_id }}
        This is you!
      {{ else if .followed }}
        You are currently following this user.
        <a class=unfollow href="/unfollow_user?username={{.profile_user.username}}">Unfollow user</a>.
      {{ else }}
        You are not yet following this user.
        <a class=follow href="/follow_user?username={{.profile_user.username}}">Follow user</a>.
      {{ end }}
      </div>
    {{ else if .endpoint eq "timeline" }}
      <div class=twitbox>
        <h3>What's on your mind {{ .user.username }}?</h3>
        <form action="/add_message" method=post>
          <p><input type=text name=text size=60><!--
          --><input type=submit value="Share">
        </form>
      </div>
    {{ end }}
  {{ end }}
  <ul class=messages>
  {{ range .messages }}
    <li><img src="http://www.gravatar.com/avatar/{{ .message_id }}?d=identicon&amp;s=48"><p>
      <strong><a href="/user_timeline?username={{.author}}">{{ .author }}</a></strong>
      {{ .text }}
      <small>&mdash; {{ .created_at }}</small>
  {{ else }}
    <li><em>There's no message so far.</em>
  {{ end }}
  </ul>
  </div>
  <div class=footer>
    MiniTwit &mdash; A Gin Application
  </div>
</div>
