'use client';

import React, { useMemo, useState } from 'react';
import styled, { useTheme } from 'styled-components';
import { useGoEnterpriseOffsets } from '@/hooks';
import { ChevronDownIcon, ChevronRightIcon, GoLogo } from '@odigos/ui-kit/icons';
import {
  FlexColumn,
  FlexRow,
  Loader,
  PageContent,
  Search,
  Tag,
  Typography,
  TypographyColor,
  TypographySize,
  TypographyVariants,
} from '@odigos/ui-kit/components';

const Header = styled(FlexColumn)`
  margin-bottom: 8px;
`;

const Meta = styled(Typography)`
  opacity: 0.8;
`;

const List = styled(FlexColumn)`
  flex: 1;
  min-height: 0;
  overflow: auto;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
`;

const RowButton = styled.button<{ $expanded: boolean }>`
  all: unset;
  box-sizing: border-box;
  display: grid;
  grid-template-columns: 24px minmax(0, 2fr) minmax(120px, 0.7fr) minmax(120px, 0.7fr);
  gap: 12px;
  align-items: center;
  width: 100%;
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme, $expanded }) => ($expanded ? theme.colors.dropdown_bg_2 : 'transparent')};

  &:hover {
    background: ${({ theme }) => theme.colors.dropdown_bg_2};
  }
`;

const ListHeader = styled.div`
  display: grid;
  grid-template-columns: 24px minmax(0, 2fr) minmax(120px, 0.7fr) minmax(120px, 0.7fr);
  gap: 12px;
  align-items: center;
  width: 100%;
  padding: 10px 16px;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  position: sticky;
  top: 0;
  background: ${({ theme }) => theme.colors.dropdown_bg};
`;

const ExpandedPanel = styled.div`
  padding: 12px 16px 16px 52px;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme }) => theme.colors.dropdown_bg};
`;

const VersionsWrap = styled(FlexRow)`
  flex-wrap: wrap;
`;

const EmptyState = styled(FlexColumn)`
  padding: 48px 16px;
  align-items: center;
`;

const Truncate = styled(Typography)`
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;

export default function Page() {
  const theme = useTheme();
  const { goEnterpriseOffsets, loading, error } = useGoEnterpriseOffsets();
  const [search, setSearch] = useState('');
  const [expandedModule, setExpandedModule] = useState<string | null>(null);

  const filteredMods = useMemo(() => {
    const mods = goEnterpriseOffsets?.mods ?? [];
    const q = search.trim().toLowerCase();
    if (!q) return mods;
    return mods.filter(
      (mod) =>
        mod.module.toLowerCase().includes(q) ||
        mod.minVersion.toLowerCase().includes(q) ||
        mod.maxVersion.toLowerCase().includes(q) ||
        mod.versions.some((v) => v.toLowerCase().includes(q)),
    );
  }, [goEnterpriseOffsets?.mods, search]);

  const timestampLabel = goEnterpriseOffsets?.timestamp
    ? new Date(goEnterpriseOffsets.timestamp).toLocaleString()
    : undefined;

  const toggleModule = (moduleName: string) => {
    setExpandedModule((current) => (current === moduleName ? null : moduleName));
  };

  return (
    <PageContent $gap={16}>
      <Header $gap={8}>
        <FlexRow $gap={12} $alignItems='center'>
          <GoLogo />
          <Typography variant={TypographyVariants.H2} size={TypographySize.XL}>
            Go Enterprise Offsets
          </Typography>
        </FlexRow>
        <Meta variant={TypographyVariants.P} size={TypographySize.S} color={TypographyColor.Secondary}>
          {error
            ? error.message
            : timestampLabel
              ? `Last updated ${timestampLabel} · ${filteredMods.length} modules`
              : loading
                ? 'Loading go offsets…'
                : `${filteredMods.length} modules`}
        </Meta>
      </Header>

      <Search data-id='go-offsets-search' value={search} onChange={setSearch} placeholder='Search modules or versions' width={360} />

      <List>
        <ListHeader>
          <span />
          <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
            Module
          </Typography>
          <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
            Min version
          </Typography>
          <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
            Max version
          </Typography>
        </ListHeader>

        {loading && (
          <EmptyState $gap={12}>
            <Loader withSpinnerOld scaleSpinnerOld={1.5} />
          </EmptyState>
        )}

        {!loading && filteredMods.length === 0 && (
          <EmptyState $gap={8}>
            <Typography variant={TypographyVariants.P} size={TypographySize.M}>
              No go offsets
            </Typography>
            <Typography variant={TypographyVariants.P} size={TypographySize.S} color={TypographyColor.Secondary}>
              {error ? error.message : 'No modules found in the odigos-go-offsets ConfigMap.'}
            </Typography>
          </EmptyState>
        )}

        {!loading &&
          filteredMods.map((mod) => {
            const expanded = expandedModule === mod.module;
            return (
              <React.Fragment key={mod.module}>
                <RowButton $expanded={expanded} onClick={() => toggleModule(mod.module)} type='button'>
                  {expanded ? <ChevronDownIcon fill={theme.colors.secondary} /> : <ChevronRightIcon fill={theme.colors.secondary} />}
                  <Truncate variant={TypographyVariants.Span} size={TypographySize.S} title={mod.module}>
                    {mod.module}
                  </Truncate>
                  <Typography variant={TypographyVariants.Span} size={TypographySize.S}>
                    {mod.minVersion || '—'}
                  </Typography>
                  <Typography variant={TypographyVariants.Span} size={TypographySize.S}>
                    {mod.maxVersion || '—'}
                  </Typography>
                </RowButton>
                {expanded && (
                  <ExpandedPanel>
                    <VersionsWrap $gap={8}>
                      {mod.versions.length === 0 ? (
                        <Typography variant={TypographyVariants.Span} size={TypographySize.S} color={TypographyColor.Secondary}>
                          No versions
                        </Typography>
                      ) : (
                        mod.versions.map((version) => <Tag key={version} label={version} />)
                      )}
                    </VersionsWrap>
                  </ExpandedPanel>
                )}
              </React.Fragment>
            );
          })}
      </List>
    </PageContent>
  );
}
