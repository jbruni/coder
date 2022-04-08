import List from "@material-ui/core/List"
import ListItem from "@material-ui/core/ListItem"
import { makeStyles } from "@material-ui/core/styles"
import React from "react"
import { NavLink } from "react-router-dom"
import { combineClasses } from "../../util/combineClasses"

export interface TabSidebarItem {
  path: string
  label: string
  hasChanges?: boolean
}

export interface TabSidebarProps {
  menuItems: TabSidebarItem[]
}

export const TabSidebar: React.FC<TabSidebarProps> = ({ menuItems }) => {
  const styles = useStyles()

  return (
    <List className={styles.menu}>
      {menuItems.map(({ hasChanges, ...tab }) => {
        return (
          <NavLink to={tab.path} key={tab.path} className={styles.link}>
            {({ isActive }) => (
              <ListItem button className={styles.menuItem} disableRipple focusRipple={false} component="li">
                <span className={combineClasses({ [styles.menuItemSpan]: true, active: isActive })}>
                  {hasChanges ? `${tab.label}*` : tab.label}
                </span>
              </ListItem>
            )}
          </NavLink>
        )
      })}
    </List>
  )
}

const useStyles = makeStyles((theme) => ({
  menu: {
    minWidth: 160,
    marginTop: theme.spacing(5),
  },

  link: {
    textDecoration: "none",
  },

  menuItem: {
    letterSpacing: -theme.spacing(0.0375),
    padding: 0,
    fontSize: 18,
    color: theme.palette.text.secondary,

    "&.MuiListItem-button:hover": {
      backgroundColor: "transparent",
    },
  },

  menuItemSpan: {
    paddingTop: theme.spacing(1.5),
    paddingBottom: theme.spacing(1),
    transition: "300ms all ease",
    position: "relative",

    "&:hover": {
      color: theme.palette.text.primary,
    },
    "&.active": {
      color: theme.palette.primary.dark,
    },

    "&.active, &:hover": {
      "&::before": {
        opacity: 1,
      },
      "&::after": {
        opacity: 1,
      },
    },
  },

  [theme.breakpoints.up("lg")]: {
    menuItemSpan: {
      "&::before": {
        content: "'{'",
        position: "absolute",
        left: -15,
        opacity: 0,
        transition: "inherit",
      },
      "&::after": {
        content: "'}'",
        position: "absolute",
        right: -15,
        opacity: 0,
        transition: "inherit",
      },
    },
  },
}))
