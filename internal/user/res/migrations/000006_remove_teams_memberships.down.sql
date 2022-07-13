CREATE TABLE IF NOT EXISTS teams (
    team_id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(32) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (user_id) REFERENCES user_profiles (user_id)
);

CREATE INDEX idx_team_tenant ON teams(tenant_id);
CREATE INDEX idx_team_name ON teams(name);

CREATE TABLE IF NOT EXISTS memberships (
    membership_id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    team_id VARCHAR(36) NOT NULL,
    role ROLE DEFAULT 'editor',
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (user_id) REFERENCES user_profiles (user_id),
    FOREIGN KEY (team_id) REFERENCES teams (team_id)
);

CREATE INDEX idx_membership_tenant ON memberships(tenant_id);
CREATE INDEX idx_memberships_user ON memberships(user_id);

-- projects is a redundant copy of project service data
CREATE TABLE projects (
    project_id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(36) NOT NULL,
    prefix VARCHAR(4) NOT NULL,
    description TEXT,
    user_id VARCHAR(36) NOT NULL,
    team_id VARCHAR(36),
    active BOOLEAN DEFAULT TRUE,
    public BOOLEAN DEFAULT FALSE,
    column_order TEXT ARRAY[10],
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (user_id) REFERENCES user_profiles (user_id)
);

CREATE INDEX idx_project_tenant ON projects(tenant_id);
CREATE INDEX idx_project_team ON projects(team_id);
CREATE INDEX idx_project_user ON projects(user_id);

ALTER TABLE invites ADD COLUMN team_id VARCHAR(36) NOT NULL DEFAULT 'n/a';
ALTER TABLE invites ADD CONSTRAINT invites_team_id_fkey FOREIGN KEY (user_id) REFERENCES user_profiles (user_id);

-- enable RLS
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE memberships ENABLE ROW LEVEL SECURITY;

-- create policies
CREATE POLICY projects_isolation_policy ON projects
    USING (tenant_id = (SELECT current_setting('app.current_tenant')));

CREATE POLICY teams_isolation_policy ON teams
    USING (tenant_id = (SELECT current_setting('app.current_tenant')));

CREATE POLICY memberships_isolation_policy ON memberships
    USING (tenant_id = (SELECT current_setting('app.current_tenant')));
