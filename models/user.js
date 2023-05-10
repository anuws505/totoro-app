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
    }
  },
  {
    timestamps: true
  }
)

module.exports = mongoose.model('User', userSchema)
