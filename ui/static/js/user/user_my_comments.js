const container = document.getElementById('commentContainer');
const template = document.getElementById('comment-template');

window.addEventListener('DOMContentLoaded', () => {
  loadComments();
});

async function loadComments() {
  try {
    const resp = await fetch('http://localhost:8080/forum/api/user/comments', {
      credentials: 'include'
    });
    if (!resp.ok) throw new Error('Failed to load comments');
    const comments = await resp.json();
    render(comments);
  } catch (err) {
    console.error(err);
    container.textContent = 'No comments found.';
  }
}

function render(comments) {
  container.innerHTML = '';
  if (!comments.length) {
    container.textContent = 'No comments found.';
    return;
  }
  comments.forEach(c => {
    const node = template.content.cloneNode(true);
    node.querySelector('.comment-post').textContent = c.post_title;
    node.querySelector('.comment-content').textContent = c.content;
    node.querySelector('.comment-time').textContent = new Date(c.created_at).toLocaleString();
    const wrapper = document.createElement('a');
    wrapper.href = `/user/post?id=${c.post_id}&highlight=${c.id}`;
    wrapper.className = 'post-link';
    wrapper.appendChild(node);
    container.appendChild(wrapper);
  });
}
