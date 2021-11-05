import { Router } from 'itty-router'
import { json, missing, withContent } from 'itty-router-extras'

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

router.get('/posts', async request => {
  var json_keys = await posts.list()
  var arr_keys = json_keys.keys

  console.log(arr_keys)

  const ret = await Promise.all(arr_keys.map(key => posts.get(key.name)))
  console.log(ret)
  // return posts as json
  return new Response(JSON.stringify(ret), {
    headers: {
      'Content-Type': 'application/json',
    },
  })
})

function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
    (
      c ^
      (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))
    ).toString(16),
  )
}

router.post('/posts', withContent, async request => {
  try {
    const { content } = request
    console.log(content)
    var body = JSON.stringify(content)
    var escaped = body.replace(/\\n/g, "\\n")
                                      .replace(/\\'/g, "\\'")
                                      .replace(/\\"/g, '\\"')
                                      .replace(/\\&/g, "\\&")
                                      .replace(/\\r/g, "\\r")
                                      .replace(/\\t/g, "\\t")
                                      .replace(/\\b/g, "\\b")
                                      .replace(/\\f/g, "\\f")
    console.log(escaped)
    await posts.put(uuidv4(), escaped)

    return new Response(escaped, {
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (err) {
    return new Response(err, { status: 500 })
  }
})

router.all('*', () => new Response('404, not found!', { status: 404 }))

addEventListener('fetch', e => {
  e.respondWith(router.handle(e.request))
})
