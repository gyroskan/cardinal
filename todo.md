# TODO

## Urgent

### Generic

- [ ] StatusInternalServerError with better description
- [ ] Proper Sql error handling
  - [ ] SqlError: Duplicate entry
  - [ ] SqlError: Foreign key constraint (invalid guild / invalid member ...)
- [ ] Rework logging

### Members

- ResetGuildMembers
  - [ ] 404 on guild not found
- ResetMember
  - [ ] InternalServerError with description and maybe a field indicating if the member was reset
  - [ ] StatusNotFound proper map
- UpdateMember
  - [ ] Better BadRequest (not sending the error alone)

### Roles

- getRoles
  - [ ] Better query
  - [ ] Reward might allow sql injection
- createRole
  - [ ] Duplicate role sql error
- deleteRole
  - [ ] NotFound error

----------

## Tests

### Roles

- [ ] getRoles invalid query params

----------

## Improvements

### Sql

- [ ] Commit & transaction for rollback

### Globale structure

- [ ] DTOs for api requests/responses
- [ ] Models for database
- [ ] Converters / mappers DTOs $\leftrightarrow$ Models
