import React, { useState, useEffect } from 'react';
import axios from 'axios';

const ContainerTable = () => {
  const [containers, setContainers] = useState([]);
  const [editingContainer, setEditingContainer] = useState(null);
  const [newContainer, setNewContainer] = useState({ address: '' });
  const [showDetails, setShowDetails] = useState(null);
  const [sortConfig, setSortConfig] = useState({
    key: 'id',
    direction: 'asc'
  });

  // Fetch all containers
  useEffect(() => {
    fetchContainers();
  }, []);

  // Сортировка контейнеров
  const sortedContainers = React.useMemo(() => {
    const sortableItems = [...containers];
    if (sortConfig.key) {
      sortableItems.sort((a, b) => {
        if (a[sortConfig.key] < b[sortConfig.key]) {
          return sortConfig.direction === 'asc' ? -1 : 1;
        }
        if (a[sortConfig.key] > b[sortConfig.key]) {
          return sortConfig.direction === 'asc' ? 1 : -1;
        }
        return 0;
      });
    }
    return sortableItems;
  }, [containers, sortConfig]);

  const handleSort = (key) => {
    let direction = 'asc';
    if (sortConfig.key === key && sortConfig.direction === 'asc') {
      direction = 'desc';
    }
    setSortConfig({ key, direction });
  };

  const fetchContainers = async () => {
    try {
      const response = await axios.get('http://localhost:8080/containers');
      setContainers(response.data);
    } catch (error) {
      console.error('Error fetching containers:', error);
    }
  };

  // Create container
  const handleCreate = async () => {
    try {
      await axios.post('http://localhost:8080/container', newContainer);
      setNewContainer({ address: '' });
      fetchContainers();
    } catch (error) {
      console.error('Error creating container:', error);
    }
  };

  // Update container
  const handleUpdate = async () => {
    try {
      await axios.put('http://localhost:8080/container', editingContainer);
      setEditingContainer(null);
      fetchContainers();
    } catch (error) {
      console.error('Error updating container:', error);
    }
  };

  // Delete container
  const handleDelete = async (id) => {
    try {
      await axios.delete(`http://localhost:8080/container/${id}`);
      fetchContainers();
    } catch (error) {
      console.error('Error deleting container:', error);
    }
  };

  // Date formatting helper
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  return (
    <div className="container">
      <h2>Container Listing</h2>
      
      {/* Add New Form */}
      <div className="add-form">
        <input
          type="text"
          placeholder="Enter IP Address"
          value={newContainer.address}
          onChange={(e) => setNewContainer({ address: e.target.value })}
        />
        <button onClick={handleCreate}>Add New</button>
      </div>

      {/* Table */}
      <table>
        <thead>
          <tr>
            <th 
              onClick={() => handleSort('id')}
              style={{ cursor: 'pointer' }}
            >
              ID {sortConfig.key === 'id' && (
                sortConfig.direction === 'asc' ? '↑' : '↓'
              )}
            </th>
            <th>IP Address</th>
            <th>Last Ping</th>
            <th>Last Success Ping</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {sortedContainers.map((container) => (
            <tr key={container.id}>
              <td>{container.id}</td>
              <td>
                {editingContainer?.id === container.id ? (
                  <input
                    value={editingContainer.address}
                    onChange={(e) => setEditingContainer({
                      ...editingContainer,
                      address: e.target.value
                    })}
                  />
                ) : (
                  container.address
                )}
              </td>
              <td>{formatDate(container.last_ping)}</td>
              <td>{formatDate(container.last_success_ping)}</td>
              <td>
                {editingContainer?.id === container.id ? (
                  <>
                    <button onClick={handleUpdate}>Save</button>
                    <button onClick={() => setEditingContainer(null)}>Cancel</button>
                  </>
                ) : (
                  <>
                    <button onClick={() => setEditingContainer(container)}>Edit</button>
                    <button onClick={() => handleDelete(container.id)}>Remove</button>
                    <button onClick={() => setShowDetails(container)}>Details</button>
                  </>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Details Modal */}
      {showDetails && (
        <div className="modal">
          <h3>Container Details</h3>
          <p>ID: {showDetails.id}</p>
          <p>IP Address: {showDetails.address}</p>
          <p>Last Ping: {formatDate(showDetails.last_ping)}</p>
          <p>Last Success Ping: {formatDate(showDetails.last_success_ping)}</p>
          <button onClick={() => setShowDetails(null)}>Close</button>
        </div>
      )}
    </div>
  );
};

export default ContainerTable;