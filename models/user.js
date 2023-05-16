const mongoose = require('mongoose')

const userSchema = mongoose.Schema(
  {
    uname: {
      type: String,
      required: true
    },
    passwd: {
      type: String,
      required: true
    },
    role: {
      type: String,
      default: 'customer'
    },
    created_date: {
      type: String
    },
    updated_date: {
      type: String
    }
  }
)

module.exports = mongoose.model('User', userSchema)
