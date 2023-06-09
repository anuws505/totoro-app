const express = require('express')
const morgan = require('morgan')
const cors = require('cors')
const mongoose =  require('mongoose')
require('dotenv').config()

const app = express()

app.use(express.json())
app.use(morgan('dev'))
app.use(cors())

mongoose.connect(
  process.env.DATABASE,
  {
    useNewUrlParser: true,
    useUnifiedTopology: true
  }
)
.then(() => {
  console.log('Database connected.')
})
.catch((error) => {
  console.error('Connection error:', error)
})

/* initial route */
app.get('/', (req, res) => {
  res.send('Hello World!')
})

const userRouter = require('./routes/user')
app.use('/user', userRouter)

const numberRouter = require('./routes/number')
app.use('/number', numberRouter)

const port = 4000
app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
