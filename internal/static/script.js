function showPreloader() {
    document.getElementById('mask').style.display = 'block';
    document.getElementById('preloader').style.display = 'block';
}

window.addEventListener('load', function () {
    document.getElementById('mask').style.display = 'none';
    document.getElementById('preloader').style.display = 'none';
});