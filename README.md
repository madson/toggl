# toggl
Go package that abstracts Toggl API v8

## Example usage

#### Instantiating a client

```go
client := toggl.NewClient(token)
```

#### Getting time entries

```go
timeEntries, err := client.GetTimeEntries(startDate, endDate)
```

#### Adding a new time entry

```go
newTimeEntry := toggl.TimeEntry{...}
timeEntry, err := client.AddTimeEntry(newTimeEntry)
```

#### Getting tags

```go
tags, err := client.GetTags(workspaceID)
```

#### Adding a new tag

```go
tag, err := client.AddTag("My new tag", workspaceID)
```
