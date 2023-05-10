const Number = require('../models/number')
// const moment = require('moment-timezone')
// const dateThai = moment().tz("Asia/Bangkok").format()

exports.getNumberAll = async (req, res) => {
  const resp = {
    code: 50000,
    message: 'unknown error'
  }

  await Number.find({})
  .then((data) => {
    resp.code = 20000
    resp.message = 'success'
    resp.data = data
    res.json(resp)
  })
  .catch((err) => {
    resp.code = 40000
    resp.message = 'fail'
    res.json(resp)
  })
}

exports.create = async (req, res) => {
  const resp = {
    code: 50000,
    message: 'unknown error'
  }
  const payload = req.body
  resp.data = payload
  // res.json(resp)

  switch(true) {
    case !payload.customer:
      resp.code = 40000
      resp.message = 'require customer'
      return res.json(resp)
      break

    case !payload.sale:
      resp.code = 40000
      resp.message = 'require sale'
      return res.json(resp)
      break

    case !payload.dateSold:
      resp.code = 40000
      resp.message = 'require date sold'
      return res.json(resp)
      break

    case !payload.grandTotal:
      resp.code = 40000
      resp.message = 'require grand total'
      return res.json(resp)
      break
  }

  if (payload.customer && payload.sale && payload.dateSold && payload.grandTotal) {
    const numberData = {}
    numberData.customer = payload.customer
    numberData.sale = payload.sale
    numberData.dateSold = payload.dateSold
    numberData.grandTotal = payload.grandTotal
    resp.data2 = numberData

    await new Number(numberData).save()
    .then((data) => {
      resp.code = 20000
      resp.message = 'success'
      resp.data = data
      res.json(resp)
    })
    .catch((err) => {
      resp.code = 40000
      resp.message = 'fail'
      res.json(resp)
    })
  }
}
