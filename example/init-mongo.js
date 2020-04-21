db.createUser(
  {
    user: 'tiger',
    pwd: 'meow',
    roles : [
      {
        role: 'readWrite',
        db: 'zoo'
      }
    ]
  }
)
