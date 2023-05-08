const express = require('express')
const morgan = require('morgan')
const cors = require('cors')
const mongoose =  require('mongoose')
require('dotenv').config()

const app = express()

app.use(express.json())
app.use(morgan('dev'))
app.use(cors())

/* initial route */
app.get('/', (req, res) => {
  res.send('Hello World!')
})

const usersRouter = require('./routes/users')
app.use('/users', usersRouter)

const port = 4000
app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
