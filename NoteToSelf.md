Rework refresh token storage
Currently refresh tokens are stored under users

Ideally, refresh tokens are their own slice in the database, each token a struct with UserID (int), Token (string), and ExpiresAt (time.Time) fields.

This way the tokens can be looked up more efficiently, simply finding the specific token in the database (O(1) lookup) and then retrieving the user that matches the token (another O(1) lookup)

Maybe some of the code can be DRY'd up by making authentication middleware or helper functions for repeat logic, it might make sense to make an internal auth package

To prevent ID conflicts it's better to use UUID or something similar. Note that the sorting would then need to be updated to use time instead of ID
