import { useState, useEffect } from 'react';

interface Group {
  ID: number;
  Name: string;
  Enabled: boolean;
  Priority: number;
}

function Groups() {
  const [groups, setGroups] = useState<Group[]>([]);

  useEffect(() => {
    fetch('/api/admin/groups')
      .then(res => res.json())
      .then(data => setGroups(data));
  }, []);

  return (
    <div>
      <h1>Groups</h1>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Enabled</th>
            <th>Priority</th>
          </tr>
        </thead>
        <tbody>
          {groups.map(group => (
            <tr key={group.ID}>
              <td>{group.ID}</td>
              <td>{group.Name}</td>
              <td>{group.Enabled ? 'Yes' : 'No'}</td>
              <td>{group.Priority}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default Groups;