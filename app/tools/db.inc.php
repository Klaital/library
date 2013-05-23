<?php

class Db {
	private static $conn = null;
	public static function getInstance() {
		if (isset(Db::$conn)) {
			return Db::$conn;
		}

		Db::$conn = new Mysqli('localhost', 'root', 'h0Shinokoe', 'library');
		return Db::$conn;
	}

    public static function addNewLookupRequest($user_id, $code, $code_type,
            $sources = array(), $location = '', $notes = '') {

        $db = Db::getInstance();

        $code = $db->real_escape_string($code);
        $code_type = $db->real_escape_string($code_type);
        $location = $db->real_escape_string($location);
        $notes = $db->real_escape_string($notes);

        $sql = "INSERT INTO lookup_requests VALUES (NULL,
            $user_id, '$code_type', '$code', NOW(), '$location', '$notes')";
        $res = $db->query($sql);

        $id = $db->insert_id;
        $sql = "INSERT INTO lookup_steps VALUES ";
        $inserts = array();
        foreach ($sources as $source) {
            $inserts[] = "(NULL, $id, '$source', NOW(), NULL, NULL)";
        }

        $sql .= implode(',', $inserts);
        $res = $db->query($sql);
    }

    public static function fetchLookupRequest($request_id) {
        $db = Db::getInstance();

        $sql = "SELECT * FROM lookup_requests
            JOIN lookup_steps ON lookup_steps.request_id = lookup_requests.request_id
            WHERE lookup_requests.request_id = " . $db->real_escape_string($request_id);
        $res = $db->query($sql);

        if (!res) {
            return false;
        }

        $data = array();
        $steps = array();

        while($row = $res->fetch_assoc()) {
            $data[] = $row;
            $steps[] = array( 'step_id' => $row['step_id'],
                            'source' => $row['source'],
                            'created' => $row['created'],
                            'started' => $row['started'],
                            'completed' => $row['completed']
                    );
        }

        if (count($data) == 0) {
            return array();
        }

        return array(
                'request_id' => $data[0]['request_id'],
                'user_id' => $data[0]['user_id'],
                'code_type' => $data[0]['code_type'],
                'code' => $data[0]['code'],
                'created' => $data[0]['created'],
                'location' => $data[0]['location'],
                'notes' => $data[0]['notes'],
                'steps' => $steps
                );
    }

    public static function setLookupRequestStepStatus($step_id, $status) {
        $db = Db::getInstance();
        $sql = 'UPDATE lookup_steps SET status = ' . $db->real_escape_string($status)
            . " WHERE step_id = $step_id";
        $res = $db->query($sql);

        return (isset($res));
    }

    // Set status to IN_PROGRESS and started = NOW()
    public static function startLookupRequestStep($step_id) {
        Db::setLookupRequestStatus($step_id, 'IN_PROGRESS');
        $db = Db::getInstance();
        $sql = "UPDATE lookup_steps SET started = NOW() WHERE step_id = $step_id";
        $res = $db->query($sql);
        return isset($res);
    }

    // Set status to SUCCESS, save the title, and set the completion time
    public static function endLookupRequestStep($step_id, $status='SUCCESS', $title='') {
        $db = Db::getInstance();
        $sql = "UPDATE lookup_steps SET
            status = '" . $db->real_escape_string($status)
            . "', title = '" . $db->real_escape_string($title)
            . "', completed = NOW()";
        $res = $db->query($sql);
        return isset($res);
    }

    // Save the title, set the completed time to NOW and set status to SUCCESS
    public static function saveLookupRequestStepTitle($step_id, $title) {

        $db = Db::getInstance();
        $sql = 'UPDATE lookup_steps SET completed=NOW(), title = \'' . $db->real_escape_string($title)
            . "WHERE step_id = $step_id";
        $res = $db->query($sql);

        Db::setLookupRequestStatus($step_id, 'SUCCESS');
        return isset($res);
    }

	public static function getOutstandingLookupRequests($code_type, $max=20) {
		$db = Db::getInstance();
		$sql = "SELECT * FROM lookup_requests, lookup_steps WHERE lookup_steps.request_id = lookup_requests.request_id AND lookup_steps.status = 'SUBMITTED' LIMIT $max";
		$res = $db->query($sql);
		if (!$res) {
			return false;
		}

		$steps = array();

		while($row = $res->fetch_assoc()) {
			$steps[] = $row;
		}

		return $steps;
	}
}

