<?php

function lookupAmazonMwsProductsApi($codes, $code_type,
        $credentials = array(
            'aws_access_key_id' => '',
            'merchant_id' => '',
            'marketplace_id_list' => array(),
            'secret_key' => '',
            )) {

    // Set up the MWS connection
    $serviceUrl = 'https://mws.amazonservices.com/Products/2011-10-01';
    $config = array(
            'ServiceURL' => $serviceUrl,
            'ProxyHost' => null,
            'ProxyPort' => -1,
            'MaxErrorRetry' => 3,
        );

    $service = new MarketplaceWebServiceProducts_Client(
            $code_type['aws_access_key_id'],
            $code_type['secret_key'],
            'Klaital\'s Library Batch Lookup',
            '1.0'
            $config);

    // Set the codes to be looked up
    $request_config = array('IdType' => $code_type, 'IdList' => $codes));
    $request->setSellerId($credentials['merchant_id']);
    $response = $service->getMatchingProductForId($request);

    try {
        $resultList = $response->getGetMatchingProductForIdResult();
        foreach ($resultList as $result) {
            $result_data = array();
            if ($result->isSetId()) {
                $result_data['id'] = $result->getId();
            }
            if ($result->isSetIdType()) {
                $result_data['idType'] = $result->getIdType();
            }
            if ($result->isSetStatus()) {
                $result_data['status'] = $result->getStatus();
            }

            if ($result->isSetProducts()) {
                $products = $result->getProducts();
                $productList = $products->getProduct();
                foreach ($productList as $product) {
                    if ($product->isSetAttributeSets()) {
                        $attributeSets = $product->getAttributeSets();
                        if ($attributeSets->isSetAny()) {
                            $nodeList->getAny();
                            $result_data['title'] = $nodeList['Title'];
                            break; // only need one product per ID
                        }
                    }
                }
            }

            // Validate the data
            if (isset($result_data['id']) && isset($result_data['idType'])
                    && isset($result_data['status']) && isset($result_data['title'])) {
                return array('SUCCESS' => $result_data);
            } else {
                return array('FAILURE' => $result_data);
            }
        }
    } catch(MarketplaceWebServiceProducts_Exception $ex) {
        $s = 'Caught Exception: ' . $ex->getMessage()
            . "\nResponse Status Code: " . $ex->getStatusCode()
            . "\nError Code: " . $ex->getErrorCode()
            . "\nError Type: " . $ex->getErrorType()
            . "\nRequest ID: " . $ex->getRequestId()
            . "\nXML: " . $ex->getXML()
            . "\nResponseHeaderMetadata: " . $ex->getResponseHeaderMetadata() . "\n";
        return array('ERROR' => $s);
    }

    return array('ERROR' => 'MWS Products API lookup not yet implemented');
}
