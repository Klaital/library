<?php
require_once('db.inc.php');

function insertBatchRequests($requests) {
    $sources = array(
        'UPC' => array('searchupc.com', 'upcdatabase.com', 'upcdatabase.org', 'Amazon MWS'),
        'ISBN' => array('Google Books', 'Amazon MWS'),
    );
    foreach($requests as $request) {
        Db::addNewLookupRequest($request['user_id'], trim($request['code']), trim($request['code_type']),
            $sources[$request['code_type']], trim($request['location']), trim($request['notes']));
    }
}

function parseTsvFile($path) {
    $f = fopen($path, 'r');
    if (!$f) {
        echo "unable to open TSV file: $path";
        return false;
    }

    $requests = array();

    while($s = fgets($f)) {
        $tokens = explode("\t", $s);
        if (count($tokens) < 4 || $tokens[1] == 'Code Type'
            //    || strlen($tokens[1]) == 0 || strlen($tokens[2])
           ) {
            // This is an invalid line
            continue;
        } else {
            $request = array(
                'location' => $tokens[0],
                'code_type' => $tokens[1],
                'notes' => $tokens[3],
                'code' => $tokens[2],
                'user_id' => 1
            );

            if (strlen(trim($request['code'])) > 0 && strlen(trim($request['code_type'])) > 0) {
                $requests[] = $request;
            }
        }
    }

    fclose($f);
    return $requests;
}


//////////////
//  MAIN    //
//////////////

echo "Parsing TSV file: '$argv[1]'\n";

$requests = parseTsvFile($argv[1]);

$count = count($requests);
echo "-Parsed $count lookup requests\n";

echo "Uploading to database... ";
insertBatchRequests($requests);
echo "done\n";

