const mongoose = require('mongoose')

const numberSchema = mongoose.Schema(
  {
    customer: {
      type: String,
      required: true
    },
    sale: {
      type: String,
      required: true
    },
    dateSold: {
      type: String,
      required: true
    },
    numbers: [
      {
        number: { type: String },
        direct: {
          price: { type: Number },
          multiple: { type: Number },
          total: { type: Number },
          reward: {
            reward: { type: Boolean },
            total: { type: Number }
          }
        },
        mix: {
          price: { type: Number },
          multiple: { type: Number },
          total: { type: Number },
          reward: {
            reward: { type: Boolean },
            total: { type: Number }
          }
        }
      }
    ],
    grandTotal: {
      type: Number,
      required: true
    },
    tax: {
      percent: { type: Number },
      total: { type: Number },
      grandTotalBeforeTax: { type: Number }
    },
    created_date: {
      type: String
    },
    updated_date: {
      type: String
    }
  }
)

module.exports = mongoose.model('Number', numberSchema)
