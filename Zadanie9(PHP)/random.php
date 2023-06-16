<?php
if ($_SERVER["REQUEST_METHOD"] == "GET") {
    $limit = !empty($_GET['limit']) ? $_GET['limit'] : 'значение не передано!';
    $count = !empty($_GET['count']) ? $_GET['count'] : 'значение не передано!';

    function randomaizer($nmax)
    {
        global $uniq;
        $uniq = rand(0, $nmax);
        return $uniq;
    }
    $tresh = array();
    $i = 0;
    while ($i != $count) {
        $uniq = randomaizer($limit);
        if (!in_array($uniq, $tresh)) {
            $tresh[] = $uniq;
            echo "<p>Число : $uniq</p>";
            $i++;
        }
    }
}
?>
