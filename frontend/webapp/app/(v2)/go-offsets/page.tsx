'use client';

import React, { useMemo, useState } from 'react';
import styled, { useTheme } from 'styled-components';
import { useMutation } from '@apollo/client/react';
import { useConfig, useGoEnterpriseOffsets } from '@/hooks';
import { UPDATE_GO_ENTERPRISE_OFFSETS } from '@/graphql';
import { GO_OFFSETS_PUBLIC_URL } from '@/utils';
import { ChevronDownIcon, ChevronRightIcon, GoLogo, RefreshIcon } from '@odigos/ui-kit/icons';
import { useNotificationStore } from '@odigos/ui-kit/store';
import { StatusType } from '@odigos/ui-kit/types';
import {
  Button,
  ButtonVariants,
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

const Toolbar = styled(FlexRow)`
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
`;

const List = styled(FlexColumn)`
  flex: 1;
  min-height: 0;
  overflow: auto;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
`;

const ModuleHeader = styled.button<{ $expanded: boolean }>`
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
  z-index: 1;
`;

const ExpandedPanel = styled(FlexColumn)`
  padding: 8px 16px 16px 52px;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme }) => theme.colors.dropdown_bg};
`;

const MinorVersionRow = styled.div`
  display: grid;
  grid-template-columns: 72px minmax(0, 1fr);
  gap: 12px;
  align-items: start;
  padding: 8px 0;
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
  const { isReadonly } = useConfig();
  const { goEnterpriseOffsets, loading, error, refetch } = useGoEnterpriseOffsets();
  const { addNotification } = useNotificationStore();
  const [search, setSearch] = useState('');
  const [expandedModules, setExpandedModules] = useState<Set<string>>(new Set());

  const [updateOffsets, { loading: mutating }] = useMutation(UPDATE_GO_ENTERPRISE_OFFSETS);
  const [fetching, setFetching] = useState(false);
  const updating = fetching || mutating;

  const filteredMods = useMemo(() => {
    const mods = goEnterpriseOffsets?.mods ?? [];
    const q = search.trim().toLowerCase();
    if (!q) return mods;
    return mods.filter(
      (mod) =>
        mod.module.toLowerCase().includes(q) ||
        mod.minVersion.toLowerCase().includes(q) ||
        mod.maxVersion.toLowerCase().includes(q) ||
        mod.minorVersions.some(
          (minor) =>
            minor.minorVersion.toLowerCase().includes(q) ||
            minor.versions.some((v) => v.toLowerCase().includes(q)),
        ),
    );
  }, [goEnterpriseOffsets?.mods, search]);

  const timestampLabel = goEnterpriseOffsets?.timestamp
    ? new Date(goEnterpriseOffsets.timestamp).toLocaleString()
    : undefined;

  const toggleModule = (moduleName: string) => {
    setExpandedModules((current) => {
      const next = new Set(current);
      if (next.has(moduleName)) {
        next.delete(moduleName);
      } else {
        next.add(moduleName);
      }
      return next;
    });
  };

  const onUpdateOffsets = async () => {
    setFetching(true);
    try {
      const response = await fetch(GO_OFFSETS_PUBLIC_URL);
      if (!response.ok) {
        throw new Error(`Failed to fetch offsets (${response.status})`);
      }
      const content = await response.text();
      await updateOffsets({ variables: { content } });
      await refetch();
      addNotification({ type: StatusType.Success, message: 'Go offsets updated from the latest public manifest.' });
    } catch (err) {
      addNotification({
        type: StatusType.Error,
        message: err instanceof Error ? err.message : 'Failed to update go offsets.',
      });
    } finally {
      setFetching(false);
    }
  };

  return (
    <PageContent $gap={16}>
      <Header $gap={8}>
        <FlexRow $gap={12} $alignItems='center' $justifyContent='space-between' $wrap='wrap'>
          <FlexRow $gap={12} $alignItems='center'>
            <GoLogo />
            <Typography variant={TypographyVariants.H2} size={TypographySize.XL}>
              Go Enterprise Offsets
            </Typography>
          </FlexRow>
          <Button
            data-id='update-go-offsets'
            label='Update offsets'
            leftIcon={RefreshIcon}
            variant={ButtonVariants.Primary}
            loading={updating}
            disabled={isReadonly || updating}
            onClick={onUpdateOffsets}
          />
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

      <Toolbar $gap={12}>
        <Search data-id='go-offsets-search' value={search} onChange={setSearch} placeholder='Search modules or versions' width={360} />
      </Toolbar>

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
            const expanded = expandedModules.has(mod.module);
            return (
              <React.Fragment key={mod.module}>
                <ModuleHeader $expanded={expanded} onClick={() => toggleModule(mod.module)} type='button'>
                  {expanded ? <ChevronDownIcon fill={theme.colors.secondary} /> : <ChevronRightIcon fill={theme.colors.secondary} />}
                  <div title={mod.module}>
                    <Truncate variant={TypographyVariants.Span} size={TypographySize.S}>
                      {mod.module}
                    </Truncate>
                  </div>
                  <Typography variant={TypographyVariants.Span} size={TypographySize.S}>
                    {mod.minVersion || '—'}
                  </Typography>
                  <Typography variant={TypographyVariants.Span} size={TypographySize.S}>
                    {mod.maxVersion || '—'}
                  </Typography>
                </ModuleHeader>
                {expanded && (
                  <ExpandedPanel $gap={0}>
                    {mod.minorVersions.length === 0 ? (
                      <Typography variant={TypographyVariants.Span} size={TypographySize.S} color={TypographyColor.Secondary}>
                        No versions
                      </Typography>
                    ) : (
                      mod.minorVersions.map((minor) => (
                        <MinorVersionRow key={minor.minorVersion}>
                          <Typography variant={TypographyVariants.Span} size={TypographySize.S} color={TypographyColor.Secondary}>
                            {minor.minorVersion}
                          </Typography>
                          <VersionsWrap $gap={8}>
                            {minor.versions.map((version) => (
                              <Tag key={version} label={version} />
                            ))}
                          </VersionsWrap>
                        </MinorVersionRow>
                      ))
                    )}
                  </ExpandedPanel>
                )}
              </React.Fragment>
            );
          })}
      </List>
    </PageContent>
  );
}
