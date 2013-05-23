<h1>Edit Item</h1>
<?php
echo $this->Form->create('Item', array('action'=>'edit'));
echo $this->Form->input('code');
echo $this->Form->input('code_type');
echo $this->Form->input('item_type');
echo $this->Form->input('location');
echo $this->Form->input('title');
echo $this->Form->input('title_translated');
echo $this->Form->input('data_source');
echo $this->Form->input('comments');
echo $this->Form->input('id', array('type' => 'hidden'));
echo $this->Form->end('Save Item');

