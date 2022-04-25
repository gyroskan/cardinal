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

### Users

- Definitely need DTOs to avoid handling password in the same struct
- banUser
  - [ ] Not found on the update

----------

## Tests

### Roles

- [ ] getRoles invalid query params

### Users

- [ ] Check if user can update is own user_access without beeing admin. (with the POST("/:username"))

----------

## Improvements

### Sql

- [ ] Commit & transaction for rollback

### Globale structure

- [ ] DTOs for api requests/responses
- [ ] Models for database
- [ ] Converters / mappers DTOs $\leftrightarrow$ Models
