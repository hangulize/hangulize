import 'semantic-ui-less/semantic.less'
import './index.css'

import * as React from 'react'
import * as ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import App from './components/App'

const rootElem = document.getElementById('root') as HTMLElement
const root = ReactDOM.createRoot(rootElem)
root.render(
  <BrowserRouter>
    <App />
  </BrowserRouter>
)
