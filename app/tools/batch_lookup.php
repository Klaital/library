<?php
require_once('XML/RPC.php');
function debug_log($msg, $file='debug.log') {
	$f = fopen($file, 'a');
	fputs($f, strftime("%FT%T " . $msg."\n"));
	fflush($f);
	fclose($f);
}

function lookupUpcSearchUpcDotCom($code, $api_key = '2F26CC56-6D00-4525-B67F-5E3A13DA57CD') {

    $json_api_url =
        "http://www.searchupc.com/handlers/upcsearch.ashx?request_type=3&access_token=$api_key&upc=$code";
    debug_log('searchupc.com querying with url ' . $json_api_url);

    $get_cmd = 'curl ' . str_replace('&', '\&', $json_api_url);
    $response = `$get_cmd`;

    if ($response) {
        $data = get_object_vars(json_decode($response));
        debug_log('searchupc.com reponse: ');
        debug_log($data);
        if (isset($data[0]->productname) && !empty($data[0]->productname)) {
            debug_log('searchupc.com found item name: ' . $data[0]->productname);
            return $data[0]->productname;
        }
    } else {
        debug_log('searchupc.com lookup failed');
    }

    return false;
}

function lookupUpcDatabaseDotOrg($code, $api_key = 'dfb7fd1143b74fa0f0eb433c50992125') {
    debug_log('upcdatabase.org querying code ' . $code);

    $json_api_url = "http://www.upcdatabase.org/api/json/$api_key/$code";
    $get_cmd = 'curl ' . $json_api_url;
    $response = `$get_cmd`;
    if ($response) {
        debug_log('upcdatabase.org response: ' . $response);
        $data = get_object_vars(json_decode($response));
        if (isset($data['itemname']) && !empty($data['itemname'])) {
            debug_log('upcdatabase.org found itemname: ' . $data['itemname']);
            return $data['itemname'];
        }
    } else {
        debug_log('upcdatabase.org failed to respond');
    }
    return false;
}

function lookupUpcDatabaseDotCom($code, $rpc_key = 'd9c8d499be82d372f77f32ca507077206a97a898') {
    // Try to look up from upcdatabase.com
    debug_log('upcdatabase.com querying code ' . $code);
    $client = new XML_RPC_Client('/xmlrpc', 'http://www.upcdatabase.com');
    $params = array( new XML_RPC_Value( array(
                    'rpc_key' => new XML_RPC_Value($rpc_key, 'string'),
                    'upc' => new XML_RPC_Value($code, 'string'),
                    ), 'struct'));
    $msg = new XML_RPC_Message('lookup', $params);
    $resp = $client->send($msg);
    if ($resp) {
        $upc_data = XML_RPC_decode($resp->value());
        debug_log('upcdatabase.com lookup success for code ' . $code);
        debug_log($upc_data);
        if (isset($upc_data['description']) && !empty($upc_data['description'])) {
            return $upc_data['description'];
        } else {
            debug_log('upcdatabase.com response: No description field found');
        }
    } else {
        debug_log('updatabase.com lookup failed: ' . $client->errstr);
    }

    return false;
}

function lookupIsbnGoogle($code) {

    // Try checking books.google.com
    $url = 'http://books.google.com/books/feeds/volumes?q=isbn:'.$code;
    debug_log('Checking books.google.com for ISBN ' . $code);
    debug_log('Using URL: ' . $url);
    $curl_cmd = 'curl ' . str_replace('&', '\&', $url);
    $response = `$curl_cmd`;

    debug_log('books.google.com response: ' . $response);
        /*
        $matches = array();
        if (!preg_match('/<dc:title>.*<\/dc:title>/i', $response, $matches)) {
            debug_log('Unable to find a title in the response. Failing search...');
            return false;
        }

        $title = substr($matches[0], 10, -11);
        */
    try {
        $feed = new SimpleXMLElement($response);
    } catch (Exception $e) {
        debug_log('books.google.com Unable to find ISBN ' . $code);
        return false;
    }

    print_r($response);
    if (!preg_match('/dc:title/', $response)) {
        debug_log('books.google.com response contained no title element');
        return false;
    }

    $title = $feed->entry->title;
    debug_log('Extracted book title: ' . $title);
    if (isset($title) && $title && strlen($title) > 0) {
        return $title;
    }

    // Unable to get a book title from books.google.com
    return false;
}

function processLookupRequest($request_id) {
    $db = Db::getInstance();
    $request = Db::fetchLookupRequest($request_id);
    foreach($request['steps'] as $step) {
        processLookupRequestStep($step);
    }
}

function processLookupRequestStep($code, $code_type, $step) {
    // Mark the step as claimed
    Db::startLookupRequestStep($step['step_id']);

    // Fetch a title based on the source designated
    $title = null;
    if ($code_type == 'UPC') {
        if ($step['source'] == 'searchupc.com') {
            $title = lookupUpcSearchUpcDotCom($code);
        } else if ($step['source'] == 'upcdatabase.com') {
            $title = lookupUpcDatabaseDotCom($code);
        } else if ($step['source'] == 'upcdatabase.org') {
            $title = lookupUpcDatabaseDotOrg($code);
        } else if ($step['source'] == 'Amazon MWS') {
            // MWS lookup actually returns an array of [ 'STATUS_CODE' => 'MESSAGE' ]
            $title = lookupAmazonMwsProductsApi($code, $code_type);
            if (isset($title['ERROR'])) {
                Db::setlookupRequestStepStatus($step['step_id'], 'Error: MWS lookup error: ' . $title['ERROR']);
                return false;
            }
            if (isset($title['FAILURE'])) {
                Db::setLookupRequestStepStatus($step['step_id'], 'FAILURE');
                return false;
            }
            if (isset($title['SUCCESS'])) {
                $title = $title['SUCCESS'];
            } else {
                // Shouldn't ever get here unless I've added more return codes from the mws lookup function
                // and forgotten to add support for them here.
                Db::setLookupRequestStepStatus($step['step_id'], 'ERROR: Unknown MWS Error code');
                return false;
            }

        } else {
            // Set error state
            Db::setLookupRequestStepStatus($step['step_id'], 'Error: unknown source');
            return false;
        }
    } else if ($code_type == 'ISBN') {
        if ($step['source'] == 'Google Books') {
            $title = lookupIsbnGoogleBooks($code);
        } else if ($step['source'] == 'Amazon MWS') {
            $title = lookupAmazonMwsProductsApi($code, $code_type);
        } else {
            // Set error state
            Db::setLookupRequestStepStatus($step['step_id'], 'Error: unknown source');
            return false;
        }
    } else {
        // Set error state
        Db::setLookupRequestStepStatus($step['step_id'], 'Error: unknown code type');
        return false;
    }

    // Save the title we looked up
    //  (this function also sets the status to COMPLETE)
    Db::saveLookupRequestStepTitle($step['step_id'], $title);
    return true;
}

// --------------
// ---- MAIN ----
// --------------


