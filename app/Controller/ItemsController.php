<?php
class ItemsController extends AppController {
    public $helpers = array('Html', 'Form', 'Paginator');
    public $components = array('Session');

    public $paginate = array(
            'limit' => 25,
            'order' => array(
                'Item.location' => 'asc'
            )
    );

    var $name = 'Items';

    public function index($location=null) {
        $conditions = null;
        if (isset($location)) {
            $conditions = array(
                    'Item.location' => array($location)
            );
        }
        $this->paginate['conditions'] = $conditions
        $data = $this->paginate('Item');
        $this->set('data', $data);

        $this->set('items', $this->paginate('Item'));
        $this->set('items', $this->Item->find('all', $conditions));
    }

    public function view($id=null) {
        $this->Item->id = $id;
        $this->set('item', $this->Item->read());
    }

    public function add() {
        if ($this->request->is('post') && !empty($this->request->data)) {
            if ($this->Item->save($this->request->data)) {
                $item = $this->Item->getLastInsertID();
                if ($item) {
                    $item = $this->Item->find('first', array(
                                'conditions' => array('Item.id' => $item),
                                'fields' => array('Item.data_source', 'Item.title', 'Item.id')
                    ));
                }
                $this->log($item, 'debug');
                $this->Session->setFlash('Your item has been saved with title <i>'
                        . $item['Item']['title']
                        . '</i> from source "'
                        . $item['Item']['data_source'] . '"'
                        . ' <a href="/items/edit/' . $item['Item']['id'] . '">Edit</a>'
                );
            } else {
                $this->Session->setFlash('Unable to add your item.');
            }
        }

    }

    public function edit($id=null) {
        $this->Item->id = $id;
        if ($this->request->is('get')) {
            $this->request->data = $this->Item->read();
        } else {
            if ($this->Item->save($this->request->data)) {
                $this->Session->setFlash('Your item has been updated.');
                $this->redirect(array('action' => 'index'));
            } else {
                $this->Session->setFlash('Unable to update your item.');
            }
        }
    }

    public function delete($id) {
        if ($this->request->is('get')) {
            throw new MethodNotAllowedException();
        }
        if ($this->Item->delete($id)) {
            $this->Session->setFlash('The item with id: ' . $id . ' has been deleted.');
            $this->redirect(array('action' => 'index'));
        }
    }

    public function search() {
        if (!empty($this->data['Item']['code'])) {
            // Search the database for the given code
            $search = $this->Item->find('first', array(
                        'conditions' => array('Item.code' => $this->data['Item']['code'])
            ));

            if (isset($search['Item']['id']) && !empty($search['Item']['id'])) {
                $this->redirect(array('action' => 'view', $search['Item']['id']));
            } else {
                $this->Session->setFlash('Unable to find any item with code ' . $this->data['Item']['code']);
            }

        }
    }

}

