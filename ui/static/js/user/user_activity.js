const container = document.getElementById('activityContainer');
const postTpl = document.getElementById('post-template');
const commentTpl = document.getElementById('comment-template');

const tabPosts = document.getElementById('tab-posts');
const tabComments = document.getElementById('tab-comments');
const tabReactions = document.getElementById('tab-reactions');
const reactionsFilter = document.getElementById('reactions-filter');

let reactionsCache = [];

window.addEventListener('DOMContentLoaded', () => {
  tabPosts.addEventListener('click', loadPosts);
  tabComments.addEventListener('click', loadComments);
  tabReactions.addEventListener('click', loadReactions);
  reactionsFilter.addEventListener('change', applyReactionFilter);
  loadPosts();
});

function setActive(btn) {
  [tabPosts, tabComments, tabReactions].forEach(b => b.classList.remove('active'));
  btn.classList.add('active');

  if (btn === tabReactions) {
    reactionsFilter.style.display = 'block';
  } else {
    reactionsFilter.style.display = 'none';
  }
}

async function loadPosts() {
  setActive(tabPosts);
  try {
    const resp = await fetch('http://localhost:8080/forum/api/user/posts', { credentials: 'include' });
    if (!resp.ok) throw new Error('failed');
    const posts = await resp.json();
    renderPosts(posts);
  } catch (e) {
    container.textContent = 'No posts found.';
  }
}

async function loadComments() {
  setActive(tabComments);
  try {
    const resp = await fetch('http://localhost:8080/forum/api/user/comments', { credentials: 'include' });
    if (!resp.ok) throw new Error('failed');
    const comments = await resp.json();
    renderComments(comments);
  } catch (e) {
    container.textContent = 'No comments found.';
  }
}

async function loadReactions() {
  setActive(tabReactions);
  try {
    const resp = await fetch('http://localhost:8080/forum/api/user/reactions', { credentials: 'include' });
    if (!resp.ok) throw new Error('failed');
    reactionsCache = await resp.json();
    applyReactionFilter();
  } catch (e) {
    container.textContent = 'No reactions found.';
  }
}

function renderPosts(posts) {
  container.innerHTML = '';
  if (!posts.length) {
    container.textContent = 'No posts found.';
    return;
  }
  posts.forEach(p => {
    const node = postTpl.content.cloneNode(true);
    node.querySelector('.post-header').textContent = p.username || 'You';
    node.querySelector('.post-title').textContent = p.title;
    node.querySelector('.post-content').textContent = p.content;
    node.querySelector('.post-time').textContent = new Date(p.created_at).toLocaleString();
    const wrapper = document.createElement('a');
    wrapper.href = `/user/post?id=${p.id}`;
    wrapper.className = 'post-link';
    wrapper.appendChild(node);
    container.appendChild(wrapper);
  });
}

function renderComments(comments) {
  container.innerHTML = '';
  if (!comments.length) {
    container.textContent = 'No comments found.';
    return;
  }
  comments.forEach(c => {
    const node = commentTpl.content.cloneNode(true);
    node.querySelector('.comment-post').textContent = c.post_title;
    node.querySelector('.comment-content').textContent = c.content;
    node.querySelector('.comment-time').textContent = new Date(c.created_at).toLocaleString();
    const wrapper = document.createElement('a');
    wrapper.href = `/user/post?id=${c.post_id}`;
    wrapper.className = 'post-link';
    wrapper.appendChild(node);
    container.appendChild(wrapper);
  });
}

function applyReactionFilter() {
  if (!reactionsCache.length) {
    container.textContent = 'No reactions found.';
    return;
  }
  const value = document.querySelector('input[name="reactFilter"]:checked').value;
  let filtered = reactionsCache;
  if (value === 'liked') {
    filtered = reactionsCache.filter(p => p.reactions.some(r => r.user_id && r.reaction_type === 1));
  } else if (value === 'disliked') {
    filtered = reactionsCache.filter(p => p.reactions.some(r => r.user_id && r.reaction_type === 2));
  }
  renderPosts(filtered);
}
