package email

const emailBodyTpl = `
<!doctype html>
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>{{ build.ProjectName }} {{ build.Status }}</title>
</head>

<body bgcolor="#f6f6f6">
    <h2>{{ build.ProjectName }} ({{ build.Branch }})</h2>
    <p><b><a href="{{APP_ADDRESS}}/project-{{ build.ProjectID }}/build-{{ build.ID }}/">Build #{{ build.ID }} has {% if build.Status == "success" %}passed{% else %}failed{% endif %}</a></b>
    </p>
    {% if build.Error %} <pre>{{ build.Error }}</pre>{% endif %}
    <p><b>Commiter:</b> {{ build.CommitterName }} ({{ build.CommitterEmail }})</p>
    <p><b>Commited:</b> {{ build.CommittedAt | time:"Mon, 02 Jan 2006 15:04:05 -0700"}}</p>
    <p><b>Message:</b></p>
    <pre>{{ build.CommitMessage }}</pre>
</body>

</html>

`
