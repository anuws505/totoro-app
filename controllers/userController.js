const User = require('../models/user')
const moment = require('moment-timezone')
// const dateThai = new moment().tz('Asia/Bangkok').format()

exports.getUserAll = async (req, res) => {
  const resp = {
    code: 50000,
    message: 'unknown error'
  }

  await User.find({})
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
  const {uname, passwd, role} = req.body
  // resp.data = {uname, passwd, role}

  switch(true) {
    case !uname:
      resp.code = 40000
      resp.message = 'require uname'
      return res.json(resp)
      break

    case !passwd:
      resp.code = 40000
      resp.message = 'require passwd'
      return res.json(resp)
      break
  }

  if (uname && passwd) {
    const userData = {}
    userData.uname = uname
    userData.passwd = passwd
    userData.created_date = new moment().tz('Asia/Bangkok').format()
    userData.updated_date = new moment().tz('Asia/Bangkok').format()
    if (role && role !== '') { userData.role = role }

    await new User(userData).save()
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

exports.getUserById = async (req, res) => {
  const resp = {
    code: 50000,
    message: 'unknown error'
  }
  const {id} = req.params
  // resp.data = {id}

  switch(true) {
    case !id:
      resp.code = 40000
      resp.message = 'require user id'
      return res.json(resp)
      break
  }

  if (id) {
    await User.findById(id)
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

exports.updateUserById = async (req, res) => {
  const resp = {
    code: 50000,
    message: 'unknown error'
  }
  const {id} = req.params
  // resp.data = {id}

  switch(true) {
    case !id:
      resp.code = 40000
      resp.message = 'require user id'
      return res.json(resp)
      break
  }

  if (id) {
    const userData = {}
    if (req.body.uname && req.body.uname !== '') { userData.uname = req.body.uname }
    if (req.body.passwd && req.body.passwd !== '') { userData.passwd = req.body.passwd }
    if (req.body.role && req.body.rold !== '') { userData.role = req.body.role }
    userData.updated_date = moment().tz('Asia/Bangkok').format()

    await User.findByIdAndUpdate(id, { $set: userData })
    .then(() => {
      User.findById(id)
      .then((data) => {
        resp.code = 20000
        resp.message = 'success'
        resp.data = data
        res.json(resp)
      })
    })
    .catch((err) => {
      resp.code = 40000
      resp.message = 'fail'
      res.json(resp)
    })
  }
}

exports.deleteUserById = async (req, res) => {
  const resp = {
    code: 50000,
    message: 'unknown error'
  }
  const {id} = req.params
  // resp.data = {id}

  switch(true) {
    case !id:
      resp.code = 40000
      resp.message = 'require user id'
      return res.json(resp)
      break
  }

  if (id) {
    await User.findByIdAndDelete(id)
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
