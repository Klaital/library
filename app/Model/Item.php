<?php
#require_once('XML/RPC.php');

class Item extends AppModel {
    var $name = 'Item';
/*
 * Validation is failing with the message: Delimiter must not be alphanumeric or backslash [CORE/Cake/Model/Model.php,
 * line 3198]
 *    var $validate = array(
        'code' => array(
            'required' => true,
            'message' => 'A code must be specified',
        ),
    );
*/

    // Whenever a new Item is saved, perform internet lookups to find information about the item.
    function beforeSave() {
        if (empty($this->data['Item']['title'])) {
            if ($this->data['Item']['code_type'] == 'UPC') {

                // Try lookups in all of the UPC lookup databases
                $upc_results = array(
                        'searchupc.com' => $this->lookupUpcSearchUpcDotCom($this->data['Item']['code']),
                        'upcdatabase.com' => false,
                        'upcdatabase.org' => false);

                if (!$upc_results['searchupc.com']) {
                   $upc_results['upcdatabase.org'] = $this->lookupUpcDatabaseDotOrg($this->data['Item']['code']);
                }

                // TODO: make these multithreaded lookups
                if (!$upc_results['upcdatabase.org']) {
                    $upc_results['upcdatabase.com'] = $this->lookupUpcDatabaseDotCom($this->data['Item']['code']);
                }

                // Check to see if any of the UPC databases had a hit.
                $hit = false;
                foreach ($upc_results as $site => $result) {
                    if ($result) {
                        $hit = $site;
                        break;
                    }
                }

                if ($hit) {
                    $this->log("Using Item title from $hit: " .$upc_results[$hit], 'debug');
                    $this->data['Item']['title'] = $upc_results[$hit];
                    $this->data['Item']['data_source'] = $hit;
                    return true;
                }

                // Failed to find a hit in any of the UPC databases
                return false;

            } else if ($this->data['Item']['code_type'] == 'ISBN') {
                $this->data['Item']['data_source'] = 'books.google.com';
                $isbn_results = array('books.google.com' => false);

                $isbn_results['books.google.com'] = $this->lookupIsbnGoogle($this->data['Item']['code']);
                // Check to see if any of the UPC databases had a hit.
                $hit = false;
                foreach ($isbn_results as $site => $result) {
                    if ($result) {
                        $hit = $site;
                        break;
                    }
                }

                if ($hit) {
                    $this->log("Using Item title from $hit: " .$isbn_results[$hit], 'debug');
                    $this->data['Item']['title'] = $isbn_results[$hit];
                    $this->data['Item']['data_source'] = $hit;
                    return true;
                }

                // Unable to look up the ISBN
                return false;
            }
        } else {
            // Existing description - if there is no data_source specified, we infer it was a manual input
            if (empty($this->data['Item']['data_source'])) {
                $this->data['Item']['data_source'] = 'manual';
            }
        }

        // default return true to avoid failing the save operation
        return true;
    }

    function lookupUpcSearchUpcDotCom($code) {
        $this->log('searchupc.com querying code ' . $code, 'debug');

        $api_key = '2F26CC56-6D00-4525-B67F-5E3A13DA57CD';
        $json_api_url =
            "http://www.searchupc.com/handlers/upcsearch.ashx?request_type=3&access_token=$api_key&upc=$code";

        App::uses('HttpSocket', 'Network/Http');
        $httpSocket = new HttpSocket();
        $response = $httpSocket->get($json_api_url);
        if ($response) {
            $data = get_object_vars(json_decode($response));
            $this->log('searchupc.com reponse: ', 'debug');
            $this->log($data, 'debug');
            if (isset($data[0]->productname) && !empty($data[0]->productname)) {
                $this->log('searchupc.com found item name: ' . $data[0]->productname, 'debug');
                return $data[0]->productname;
            }
        }

        return false;
    }

    function lookupUpcDatabaseDotOrg($code) {
        $this->log('upcdatabase.org querying code ' . $code, 'debug');

        $api_key = 'dfb7fd1143b74fa0f0eb433c50992125';
        $json_api_url = "http://www.upcdatabase.org/api/json/$api_key/$code";

        App::uses('HttpSocket', 'Network/Http');

        $httpSocket = new HttpSocket();
        $response = $httpSocket->get($json_api_url);
        if ($response) {
            $this->log('upcdatabase.org response: ' . $response, 'debug');
            $data = get_object_vars(json_decode($response));
            if (isset($data['itemname']) && !empty($data['itemname'])) {
                $this->log('upcdatabase.org found itemname: ' . $data['itemname'], 'debug');
                return $data['itemname'];
            }
        } else {
            $this->log('upcdatabase.org failed to respond', 'debug');
        }
        return false;
    }

    function lookupUpcDatabaseDotCom($code) {
        // Try to look up from upcdatabase.com
#        $this->log('upcdatabase.com querying code ' . $code, 'debug');
#        $rpc_key = 'd9c8d499be82d372f77f32ca507077206a97a898';
#        $client = new XML_RPC_Client('/xmlrpc', 'http://www.upcdatabase.com');
#        $params = array( new XML_RPC_Value( array(
#                        'rpc_key' => new XML_RPC_Value($rpc_key, 'string'),
#                        'upc' => new XML_RPC_Value($code, 'string'),
#                        ), 'struct'));
#        $msg = new XML_RPC_Message('lookup', $params);
#        $resp = $client->send($msg);
#        if ($resp) {
#            $upc_data = XML_RPC_decode($resp->value());
#            $this->log('upcdatabase.com lookup success for code ' . $code, 'debug');
#            $this->log($upc_data, 'debug');
#            if (isset($upc_data['description']) && !empty($upc_data['description'])) {
#                return $upc_data['description'];
#            } else {
#                $this->log('upcdatabase.com response: No description field found', 'debug');
#            }
#        } else {
#            $this->log('updatabase.com lookup failed: ' . $client->errstr, 'debug');
#        }

        return false;
    }

    function lookupIsbnGoogle($code) {

        // Try checking books.google.com
        App::uses('HttpSocket', 'Network/Http');
//        App::import('Utility', 'Xml');
        $httpSocket = new HttpSocket();
        $url = 'http://books.google.com/books/feeds/volumes?q=isbn:'.$code;
        $this->log('Checking books.google.com for ISBN ' . $code, 'debug');
        $this->log('Using URL: ' . $url, 'debug');
        $response = $httpSocket->get($url);

        $this->log('books.google.com response: ' . $response, 'debug');
//        $feed = Xml::build($response);
//        $this->log('XML: ' . $feed, 'debug');
        // Xml::build($response) is crashing, so let's just do a simple parse for now
        $matches = array();
        if (!preg_match('/<dc:title>.*<\/dc:title>/i', $response, $matches)) {
            $this->log('Unable to find a title in the response. Failing search...', 'debug');
            return false;
        }

        $title = substr($matches[0], 10, -11);
        $this->log('Extracted book title: ' . $title, 'debug');
        if (strlen($title) > 0) {
            return $title;
        }

        // Unable to get a book title from boosk.google.com
        return false;
    }
}

