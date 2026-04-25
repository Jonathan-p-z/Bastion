(function () {
  var cards   = Array.from(document.querySelectorAll('.gcard[data-name]'));
  var total   = cards.length;
  var active  = document.querySelectorAll('.gcard-active').length;
  var countEl = document.getElementById('guilds-count');
  var emptyEl = document.getElementById('guilds-empty');
  var input   = document.getElementById('guild-search');

  if (countEl) countEl.textContent = active + ' avec Bastion · ' + total + ' au total';

  function filter(q) {
    q = q.toLowerCase().trim();
    var visible = 0;
    cards.forEach(function (c) {
      var name  = (c.getAttribute('data-name') || '').toLowerCase();
      var match = !q || name.indexOf(q) !== -1;
      c.style.display = match ? 'flex' : 'none';
      if (match) visible++;
    });
    if (emptyEl) emptyEl.style.display = visible === 0 ? '' : 'none';
  }

  if (input) {
    input.addEventListener('input',  function () { filter(this.value); });
    input.addEventListener('search', function () { filter(this.value); });
  }
})();
