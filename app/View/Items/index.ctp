<h1>Items Registered</h1>
<?php echo $this->Html->link('Add Item', array('controller' => 'items', 'action' => 'add')); ?><br />
<?php echo $this->Html->link('Find Item', array('controller' => 'items', 'action' => 'search')); ?>

<table>
    <tr>
        <th>ID</th>
        <th>Code</th>
        <th>Code Type</th>
        <th>Location</th>
        <th>Description</th>
        <th>Data Source</th>
        <th>Comments</th>
        <th>Added to Db</th>
        <th>Actions</th>
    </tr>
    <?php foreach ($items as $item): ?>
    <tr>
        <td><?php echo $this->Html->link($item['Item']['id'],
                array('controller' => 'items', 'action' => 'edit', $item['Item']['id'])); ?></td>
        <td><?php echo $this->Html->link($item['Item']['code'],
                array('controller' => 'items', 'action' => 'view', $item['Item']['id'])); ?></td>
        <td><?php echo $item['Item']['code_type']; ?></td>
        <td><?php echo $item['Item']['location']; ?></td>
        <td><?php echo $item['Item']['title']; ?></td>
        <td><?php echo $item['Item']['data_source']; ?></td>
        <td><?php echo $item['Item']['comments']; ?></td>
        <td><?php echo $item['Item']['created']; ?></td>
        <td>
            <?php echo $this->Html->link('Edit', array(
                        'controller' => 'items', 'action' => 'edit', $item['Item']['id']));
                echo ' ';
                echo $this->Form->postLink('Delete',
                        array('action' => 'delete', $item['Item']['id']),
                        array('confirm' => 'Are you sure?'));
            ?>
        </td>
    </tr>
    <?php endforeach; ?>
</table>

