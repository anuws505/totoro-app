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
    grandTotal: {
      type: Number,
      required: true
    },
    tax: {
      percent: { type: Number },
      total: { type: Number },
      grandTotalBeforetax: { type: Number }
    }
  },
  {
    timestamps: true
  }
)

module.exports = mongoose.model('Number', numberSchema)
