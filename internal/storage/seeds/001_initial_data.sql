-- Seed collections
INSERT OR IGNORE INTO collections (id, name)
VALUES
    ('jsonplaceholder', 'JSONPlaceholder'),
    ('httpbin', 'HTTPBin'),
    ('local', 'Local Development');

-- Seed requests
INSERT OR IGNORE INTO requests (
    id,
    collection_id,
    name,
    method,
    url,
    body
)
VALUES
    -- JSONPlaceholder
    (
        'jsonplaceholder-posts',
        'jsonplaceholder',
        'Get all posts',
        'GET',
        'https://jsonplaceholder.typicode.com/posts',
        ''
    ),
    (
        'jsonplaceholder-post-1',
        'jsonplaceholder',
        'Get post by ID',
        'GET',
        'https://jsonplaceholder.typicode.com/posts/1',
        ''
    ),
    (
        'jsonplaceholder-create-post',
        'jsonplaceholder',
        'Create post',
        'POST',
        'https://jsonplaceholder.typicode.com/posts',
        '{"title":"Hello World","body":"Created with net-cli","userId":1}'
    ),
    (
        'jsonplaceholder-update-post',
        'jsonplaceholder',
        'Update post',
        'PUT',
        'https://jsonplaceholder.typicode.com/posts/1',
        '{"id":1,"title":"Updated title","body":"Updated with net-cli","userId":1}'
    ),
    (
        'jsonplaceholder-delete-post',
        'jsonplaceholder',
        'Delete post',
        'DELETE',
        'https://jsonplaceholder.typicode.com/posts/1',
        ''
    ),

    -- HTTPBin
    (
        'httpbin-get',
        'httpbin',
        'Get request details',
        'GET',
        'https://httpbin.org/get',
        ''
    ),
    (
        'httpbin-post',
        'httpbin',
        'Post JSON',
        'POST',
        'https://httpbin.org/post',
        '{"message":"Hello from net-cli"}'
    ),
    (
        'httpbin-status',
        'httpbin',
        'Test status code',
        'GET',
        'https://httpbin.org/status/418',
        ''
    );