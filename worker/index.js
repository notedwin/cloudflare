import { Router } from 'itty-router'

// Create a new router
const router = Router()

/*
Our index route, a simple hello world.
*/
router.get('/', () => {
  return new Response(
    'Hello, world! This is the root page of your Worker template.',
  )
})

/*
This route demonstrates path parameters, allowing you to extract fragments from the request
URL.
Try visit /example/hello and see the response.
*/
router.get('/example/:text', ({ params }) => {
  // Decode text like "Hello%20world" into "Hello world"
  let input = decodeURIComponent(params.text)

  // Construct a buffer from our input
  let buffer = Buffer.from(input, 'utf8')

  // Serialise the buffer into a base64 string
  let base64 = buffer.toString('base64')

  // Return the HTML with the string to the client
  return new Response(`<p>Base64 encoding: <code>${base64}</code></p>`, {
    headers: {
      'Content-Type': 'text/html',
    },
  })
})


router.get('/posts', async request => {

  const value = await posts.get('first')
  if (value === null) {
    return new Response('Value not found', { status: 404 })
  }


  return new Response(value, {
    headers: { 'Content-Type': 'application/json' },
  })
})

router.all('*', () => new Response('404, not found!', { status: 404 }))

addEventListener('fetch', e => {
  e.respondWith(router.handle(e.request))
})
