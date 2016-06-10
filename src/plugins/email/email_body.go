package email

const emailBodyTpl = `
<!doctype html>
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>{{ build.ProjectName }} {{ build.Status }}</title>
</head>

<body bgcolor="#f6f6f6">

{% if build.Status == "success" %}
<b>Build #{{ build.ID }} has passed</b>
{% else %}
<b>Build #{{ build.ID }} has failed</b>
{% endif %}

<p><b>{{ build.ProjectName }}</b> (<i>{{ build.Branch }}</i>)</p>
{% if build.Error %}<pre>{{ build.Error }}</pre>{% endif %}
<p>{{ build.CommitterName }} ({{ build.CommitterEmail }}) at {{ build.CommittedAt | time:"Mon, 02 Jan 2006 15:04:05 -0700" }}</p>

<pre>{{ build.CommitMessage }}</pre>

sha: {% if APP_ADDRESS %}<a href="{{APP_ADDRESS}}/project-{{ build.ProjectID }}/build-{{ build.ID }}/">{{ build.CommitSHA }}</a>{% else %}{{ build.CommitSHA }}{% endif %}

</body>

</html>
`
